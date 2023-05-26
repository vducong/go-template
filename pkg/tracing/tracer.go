package tracing

import (
	"context"
	"fmt"
	"promotion/configs"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	trace.Tracer
}

func (t *Tracer) CustomSpan(ctx context.Context) (context.Context, trace.Span) {
	return t.Start(ctx, GetCallerInfo(3).Function)
}

func Init(cfg *configs.Config) (*Tracer, error) {
	exporter, err := texporter.New(texporter.WithProjectID(cfg.GCP.ProjectID))
	if err != nil {
		return nil, fmt.Errorf("Failed to create new exporter: %w", err)
	}

	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Server.Name),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new resource: %w", err)
	}

	sampler := getSampler(cfg.GetServerEnv())
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)
	defer traceProvider.ForceFlush(ctx)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			gcppropagator.CloudTraceFormatPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return &Tracer{traceProvider.Tracer(cfg.Server.Name)}, nil
}

func getSampler(env configs.ServerEnv) sdktrace.Sampler {
	if env == configs.ServerEnvDevelopment {
		return sdktrace.ParentBased(
			sdktrace.AlwaysSample(),
			sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
		)
	}
	return sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(0.01),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
	)
}
