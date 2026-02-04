package trc

import (
	"context"
	"fmt"
)

type Config struct {
	ServiceName string
	Kind        Kind
	Exporter    ExporterConfig
	Provider    ProviderConfig
}

type ExporterConfig struct {
	Kind     ExporterKind
	Endpoint string
}

type ProviderConfig struct {
	SamplerKind SamplerKind
	SampleRate  float64
}

type TracerProvider struct {
	StopFn func(context.Context) error
}

func Setup(ctx context.Context, configs *Config) (*TracerProvider, error) {
	if configs.Kind == KindOtel {
		return setupOtel(ctx, configs)
	}
	return nil, fmt.Errorf("invalid tracing kind=%s", configs.Kind)
}
