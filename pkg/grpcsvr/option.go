package grpcsvr

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	pbheath "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Option func(s *Server)

func WithConfig(cfg *Config) Option {
	return func(s *Server) {
		s.cfg = cfg
	}
}

func WithGrpcOptions(opts []grpc.ServerOption) Option {
	return func(s *Server) {
		s.Server = grpc.NewServer(opts...)
	}
}

func WithGrpcServices(services []Service) Option {
	return func(s *Server) {
		healthSvr := health.NewServer()
		for i := range services {
			s.Server.RegisterService(services[i].Desc(), services[i].Implementation())
			healthSvr.SetServingStatus(services[i].Desc().ServiceName, pbheath.HealthCheckResponse_SERVING)
		}
		pbheath.RegisterHealthServer(s.Server, healthSvr)

		reflection.Register(s.Server)
	}
}
