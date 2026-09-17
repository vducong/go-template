package interceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"google.golang.org/grpc/stats"
)

// TracingStatsHandler skips gRPC health checks, the gRPC match for keeping /health out of HTTP tracing.
// The filter drops their RPC metrics too.
func TracingStatsHandler() stats.Handler {
	return otelgrpc.NewServerHandler(otelgrpc.WithFilter(filters.Not(filters.HealthCheck())))
}
