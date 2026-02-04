package grpc

import (
	cfg "gotemplate/internal/config"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/grpc/interceptor"
	grpcsvc "gotemplate/internal/interface/grpc/service"
	intsvc "gotemplate/internal/service"
	"gotemplate/pkg/grpcsvr"

	"google.golang.org/grpc"
)

func New(
	configs *cfg.GRPCConfig, infrastructure *infra.Infrastructure,
	internalServices *intsvc.Services, grpcServices *grpcsvc.Services,
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
