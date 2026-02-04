package mtr

import (
	"context"
	"fmt"
	"time"
)

type Config struct {
	ServiceName string
	Kind        Kind
	Readers     []*ReaderConfig
}

type ReaderConfig struct {
	ExporterKind     ExporterKind
	ExporterEndpoint string
	Interval         time.Duration
}

type MeterProvider struct {
	StopFn     func(context.Context) error
	PromClient *PromClient
}

func Setup(ctx context.Context, configs *Config) (*MeterProvider, error) {
	if configs.Kind == KindPrometheus {
		promClient := setupProm()
		return &MeterProvider{
			PromClient: promClient,
		}, nil
	}

	if configs.Kind == KindOtel {
		return setupOtel(ctx, configs)
	}
	return nil, fmt.Errorf("invalid metric kind=%s", configs.Kind)
}
