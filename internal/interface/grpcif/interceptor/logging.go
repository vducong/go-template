package interceptor

import (
	"context"
	"gotemplate/pkg/lg"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func NewLoggingInterceptor(log lg.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := []lg.Field{
			lg.Str("protocol", "grpc"),
			lg.Str("method", info.FullMethod),
			lg.Dur("duration", duration),
		}

		if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
			spanCtx := span.SpanContext()
			fields = append(fields,
				lg.Str("trace_id", spanCtx.TraceID().String()),
				lg.Str("span_id", spanCtx.SpanID().String()),
			)
		}

		if err != nil {
			fields = append(fields, lg.Err(err))
			log.CtxError(ctx, "grpc request", fields...)
		} else {
			log.CtxInfo(ctx, "grpc request", fields...)
		}

		return resp, err
	}
}
