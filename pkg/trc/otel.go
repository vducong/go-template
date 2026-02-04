package trc

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	otelres "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

func setupOtel(ctx context.Context, configs *Config) (tracerProvider *TracerProvider, err error) {
	resources, err := otelres.Merge(
		otelres.Default(),
		otelres.NewWithAttributes(
			"",
			semconv.ServiceName(configs.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	exporter, err := initExporter(ctx, configs)
	if err != nil {
		return nil, err
	}

	stopFn, err := initTracerProvider(configs, resources, exporter)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer provider: %w", err)
	}

	return &TracerProvider{
		StopFn: stopFn,
	}, nil
}

func initExporter(ctx context.Context, configs *Config) (exporter sdktrace.SpanExporter, err error) {
	switch configs.Exporter.Kind {
	case ExporterKindStdout:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())

	case ExporterKindZipkin:
		exporter, err = zipkin.New(configs.Exporter.Endpoint)

	case ExporterKindGRPC:
		exporter, err = otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpointURL(configs.Exporter.Endpoint),
		)

	case ExporterKindHTTP:
		exporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpointURL(configs.Exporter.Endpoint),
		)

	default:
		return nil, fmt.Errorf("invalid trace exporter kind=%s", configs.Exporter.Kind)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"kind=%s | endpoint=%s | failed to init exporter: %w",
			configs.Exporter.Kind, configs.Exporter.Endpoint, err,
		)
	}
	return exporter, nil
}

func initTracerProvider(
	configs *Config,
	resource *otelres.Resource,
	exporter sdktrace.SpanExporter,
) (func(context.Context) error, error) {
	var sampler sdktrace.Sampler
	switch configs.Provider.SamplerKind {
	case SamplerKindAlways:
		sampler = sdktrace.AlwaysSample()
	case SamplerKindNever:
		sampler = sdktrace.NeverSample()
	case SamplerKindRatio:
		sampler = sdktrace.TraceIDRatioBased(configs.Provider.SampleRate)
	default:
		return nil, fmt.Errorf("invalid sampler kind=%s", configs.Provider.SamplerKind)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tracerProvider.Shutdown, nil
}
