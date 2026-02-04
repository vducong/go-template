package mtr

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

func setupOtel(ctx context.Context, configs *Config) (meterProvider *MeterProvider, err error) {
	resources, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceName(configs.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	meterProvider = &MeterProvider{}
	readers := []sdkmetric.Reader{}
	for i := range configs.Readers {
		reader, errReader := meterProvider.initReader(ctx, configs.Readers[i])
		if errReader != nil {
			return nil, fmt.Errorf("failed to init reader: %w", errReader)
		}
		readers = append(readers, reader)
	}

	stopFn, err := initMeterProvider(resources, readers)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter provider: %w", err)
	}
	return &MeterProvider{
		StopFn:     stopFn,
		PromClient: setupProm(),
	}, nil
}

func (o *MeterProvider) initReader(ctx context.Context, configs *ReaderConfig) (reader sdkmetric.Reader, err error) {
	switch configs.ExporterKind {
	case ExporterKindOtelPrometheus:
		o.PromClient = setupProm()
		reader, err = prometheus.New(
			prometheus.WithRegisterer(o.PromClient.Registry),
		)
		if err != nil {
			return nil, fmt.Errorf("kind=%s | failed to init exporter: %w", ExporterKindPrometheus, err)
		}
		return reader, nil

	case ExporterKindOtelGRPC:
		exporter, err := otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpointURL(configs.ExporterEndpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("kind=%s | endpoint=%s | failed to init exporter: %w",
				ExporterKindOtelGRPC, configs.ExporterEndpoint, err)
		}
		return sdkmetric.NewPeriodicReader(
			exporter,
			sdkmetric.WithInterval(configs.Interval),
		), nil

	case ExporterKindOtelHTTP:
		exporter, err := otlpmetrichttp.New(
			ctx,
			otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithEndpointURL(configs.ExporterEndpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("kind=%s | endpoint=%s | failed to init exporter: %w",
				ExporterKindOtelHTTP, configs.ExporterEndpoint, err)
		}
		return sdkmetric.NewPeriodicReader(
			exporter,
			sdkmetric.WithInterval(configs.Interval),
		), nil

	default:
		return nil, fmt.Errorf("invalid metric reader kind=%s", configs.ExporterKind)
	}
}

func initMeterProvider(resources *resource.Resource, readers []sdkmetric.Reader) (func(context.Context) error, error) {
	opts := []sdkmetric.Option{
		sdkmetric.WithResource(resources),
	}
	for i := range readers {
		opts = append(opts, sdkmetric.WithReader(readers[i]))
	}
	meterProvider := sdkmetric.NewMeterProvider(opts...)

	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}
