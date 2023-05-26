package server

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"promotion/configs"
	"promotion/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const cancelContextTimeout = 5 * time.Second

type Server struct {
	cfg        *configs.Config
	log        *logger.Logger
	httpServer *http.Server
}

func New(
	cfg *configs.Config, log *logger.Logger, engine *gin.Engine,
) *Server {
	s := &Server{
		cfg: cfg,
		log: log,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: engine,
		},
	}

	s.Start()
	return s
}

func (s *Server) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start HTTP server.
	go func() {
		s.serve()
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	s.log.Info("http: shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), cancelContextTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Fatalf("http: server shutdown failed: %v", err)
	}

	s.log.Info("http: server existed")
}

func (s *Server) serve() {
	s.log.Infof("http: starting server at %d", s.cfg.Server.Port)
	if err := s.httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.log.Info("http: server shutdown complete: %v", err)
		} else {
			s.log.Errorf("http: server closed unexpect: %v", err)
		}
	}
}
