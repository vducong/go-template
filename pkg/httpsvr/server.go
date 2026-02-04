package httpsvr

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type Server struct {
	server *http.Server
	ready  chan struct{}
}

func New(options ...Option) *Server {
	server := &Server{
		server: &http.Server{},
		ready:  make(chan struct{}),
	}

	for _, option := range options {
		option(server)
	}
	return server
}

func (s *Server) Start() error {
	addr := s.server.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	close(s.ready)

	if err := s.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http server: %w", err)
	}
	return nil
}

func (s *Server) WaitForReady() {
	<-s.ready
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
