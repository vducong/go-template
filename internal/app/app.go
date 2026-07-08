package app

import (
	"context"
	"fmt"
	"gotemplate/internal/apperr"
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/grpcif"
	grpcsvc "gotemplate/internal/interface/grpcif/service"
	"gotemplate/internal/interface/httpif"
	"gotemplate/internal/interface/httpif/handler"
	"gotemplate/internal/interface/httpif/middleware"
	"gotemplate/internal/repository"
	intsvc "gotemplate/internal/service"
	"gotemplate/pkg/auth"
	"gotemplate/pkg/grpcsvr"
	"gotemplate/pkg/httpsvr"
	"gotemplate/pkg/lg"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg   *cfg.Config
	infra *infra.Infrastructure

	services  *intsvc.Services
	serveDone chan struct{}

	httpServer *httpsvr.Server
	grpcServer *grpcsvr.Server

	ready chan struct{}
}

func Init(configs *cfg.Config) (*App, error) {
	app := &App{
		cfg:   configs,
		ready: make(chan struct{}),
	}
	if err := apperr.Setup(configs.App.ServiceCode, configs.App.MessageKeyPrefix); err != nil {
		return nil, fmt.Errorf("setup apperr: %w", err)
	}

	infrastructure, err := infra.Setup(configs)
	if err != nil {
		return nil, fmt.Errorf("failed to setup infra: %w", err)
	}
	app.infra = infrastructure

	repos := repository.Init(configs, infrastructure)
	services, err := intsvc.Setup(configs, infrastructure, repos)
	if err != nil {
		return nil, fmt.Errorf("failed to setup services: %w", err)
	}
	app.services = services

	if configs.HTTP.Port != "" {
		jwtCfg, err := auth.JWTConfigFromAuth(
			configs.Auth.JWTSecret,
			configs.Auth.JWTLeeway,
			configs.Auth.JWTIssuer,
			configs.Auth.JWTAudience,
		)
		if err != nil {
			return nil, fmt.Errorf("setup auth: %w", err)
		}

		authenticator, err := middleware.NewAuthenticator(
			infrastructure.Logger,
			infrastructure.ResponseWriter,
			jwtCfg,
			configs.Auth.APIKeys,
		)
		if err != nil {
			return nil, fmt.Errorf("setup authenticator: %w", err)
		}

		app.httpServer = httpif.New(
			configs, app.infra, handler.Setup(configs, infrastructure, services), authenticator,
		)
	}
	if configs.GRPC.Port != "" {
		grpcServices := grpcsvc.Setup(services)
		app.grpcServer = grpcif.New(&configs.GRPC, app.infra, grpcServices)
	}

	return app, nil
}

func (a *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if a.httpServer != nil || a.grpcServer != nil {
		errChan := make(chan error, 1)
		a.serveDone = make(chan struct{})
		go func() {
			defer close(a.serveDone)
			errChan <- a.Serve()
		}()

		select {
		case s := <-quit:
			a.infra.Logger.Info("received shutdown signal", lg.Str("signal", s.String()))
		case err := <-errChan:
			if err != nil {
				a.infra.Logger.Error("failed to serve app", lg.Err(err))
			}
		}
	} else {
		a.infra.Logger.Info("no servers configured, waiting for shutdown signal")
		s := <-quit
		a.infra.Logger.Info("received shutdown signal", lg.Str("signal", s.String()))
	}

	a.Stop()
	a.infra.Logger.Info("app shutdown complete")
}

func (a *App) Serve() error {
	var errg errgroup.Group
	if a.httpServer != nil {
		errg.Go(a.httpServer.Start)
	}
	if a.grpcServer != nil {
		errg.Go(a.grpcServer.Start)
	}

	go func() {
		if a.httpServer != nil {
			a.httpServer.WaitForReady()
			a.infra.Logger.Info("http server ready", lg.Str("port", a.cfg.HTTP.Port))
		}
		if a.grpcServer != nil {
			a.grpcServer.WaitForReady()
			a.infra.Logger.Info("grpc server ready", lg.Str("port", a.cfg.GRPC.Port))
		}
		close(a.ready)
	}()

	return errg.Wait()
}

func (a *App) waitForServe() {
	if a.serveDone == nil {
		return
	}
	select {
	case <-a.serveDone:
	case <-time.After(a.cfg.App.ShutdownTimeout):
		a.infra.Logger.Warn("background workers did not exit before timeout")
	}
}

func (a *App) Stop() {
	ctx := context.Background()

	if a.httpServer != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, a.cfg.HTTP.ShutdownTimeout)
		defer shutdownCancel()
		if err := a.httpServer.Stop(shutdownCtx); err != nil {
			a.infra.Logger.Error("failed to shutdown http server", lg.Err(err))
		}
		a.infra.Logger.Info("http server shutdown complete")
	}

	if a.grpcServer != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, a.cfg.App.ShutdownTimeout)
		defer shutdownCancel()
		a.grpcServer.Stop(shutdownCtx)
		a.infra.Logger.Info("grpc server shutdown complete")
	}

	a.waitForServe()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, a.cfg.App.ShutdownTimeout)
	defer shutdownCancel()

	if err := a.services.Stop(shutdownCtx); err != nil {
		a.infra.Logger.Error("failed to stop services", lg.Err(err))
	}
	a.infra.Logger.Info("services shutdown complete")

	// Flush tracing/metrics before closing instrumented clients (mongo, redis, gRPC).
	if err := a.infra.Close(shutdownCtx); err != nil {
		a.infra.Logger.Error("failed to close infra", lg.Err(err))
	}
	a.infra.Logger.Info("infra shutdown complete")
}
