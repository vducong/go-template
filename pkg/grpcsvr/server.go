package grpcsvr

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Port string
}

type Service interface {
	Desc() *grpc.ServiceDesc
	Implementation() any
}

type Server struct {
	Server *grpc.Server
	cfg    *Config
	ready  chan struct{}
}

func New(options ...Option) *Server {
	server := &Server{
		ready: make(chan struct{}),
	}

	for _, option := range options {
		option(server)
	}
	return server
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.cfg.Port)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}

	close(s.ready)

	if err := s.Server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve grpc server: %w", err)
	}
	return nil
}

func (s *Server) WaitForReady() {
	<-s.ready
}

func (s *Server) Stop(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		s.Server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		s.Server.Stop()
	}
}
