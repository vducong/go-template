package grpcif

import (
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/grpcif/interceptor"
	grpcsvc "gotemplate/internal/interface/grpcif/service"
	"gotemplate/pkg/grpcsvr"

	"google.golang.org/grpc"
)

func New(
	configs *cfg.GRPCConfig, infrastructure *infra.Infrastructure, grpcServices *grpcsvc.Services,
) *grpcsvr.Server {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.NewRecoveryInterceptor(infrastructure.Logger),
			interceptor.NewLoggingInterceptor(infrastructure.Logger),
		),
	}

	if infrastructure.Tracing != nil {
		opts = append(opts, grpc.StatsHandler(interceptor.TracingStatsHandler()))
	}

	return grpcsvr.New(
		grpcsvr.WithConfig(&grpcsvr.Config{
			Port: configs.Port,
		}),
		grpcsvr.WithGrpcOptions(opts),
		grpcsvr.WithGrpcServices(grpcServices.ToSlice()),
	)
}
