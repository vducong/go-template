package app

import (
	"context"
	"fmt"
	cfg "gotemplate/internal/config"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/grpc"
	"gotemplate/internal/interface/grpc/service"
	"gotemplate/internal/interface/http"
	"gotemplate/internal/interface/http/handler"
	intsvc "gotemplate/internal/service"
	"gotemplate/pkg/grpcsvr"
	"gotemplate/pkg/httpsvr"
	"gotemplate/pkg/lg"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg        *cfg.Config
	infra      *infra.Infrastructure
	httpServer *httpsvr.Server
	grpcServer *grpcsvr.Server

	ready chan struct{}
}

func Init(configs *cfg.Config) (*App, error) {
	app := &App{
		cfg:   configs,
		ready: make(chan struct{}),
	}

	infrastructure, err := infra.Setup(configs)
	if err != nil {
		return nil, fmt.Errorf("failed to setup infra: %w", err)
	}
	app.infra = infrastructure

	runtime.GOMAXPROCS(configs.App.MaxProcs)

	internalServices, err := intsvc.Setup(configs)
	if err != nil {
		return nil, fmt.Errorf("failed to setup internal services: %w", err)
	}

	if configs.HTTP.Port != "" {
		httpHandlers := handler.Setup(&handler.Config{})
		app.httpServer = http.New(configs, app.infra, httpHandlers)
	}
	if configs.GRPC.Port != "" {
		grpcServices := service.Setup(internalServices)
		app.grpcServer = grpc.New(&configs.GRPC, app.infra, internalServices, grpcServices)
	}

	return app, nil
}

func (a *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if a.httpServer != nil || a.grpcServer != nil {
		errChan := make(chan error, 1)
		go func() {
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

func (a *App) Stop() {
	if a.httpServer != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.cfg.HTTP.ShutdownTimeout)
		defer shutdownCancel()
		if err := a.httpServer.Stop(shutdownCtx); err != nil {
			a.infra.Logger.Error("failed to shutdown http server", lg.Err(err))
		}
		a.infra.Logger.Info("http server shutdown complete")
	}

	if a.grpcServer != nil {
		a.grpcServer.Stop()
		a.infra.Logger.Info("grpc server shutdown complete")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer shutdownCancel()
	if err := a.infra.Close(shutdownCtx); err != nil {
		a.infra.Logger.Error("failed to close infra", lg.Err(err))
	}
	a.infra.Logger.Info("infra shutdown complete")
}
