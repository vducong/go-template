package interceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/stats"
)

func TracingStatsHandler() stats.Handler {
	return otelgrpc.NewServerHandler()
}
