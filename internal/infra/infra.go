package infra

import (
	"context"
	"fmt"
	cfg "gotemplate/internal/config"
	"gotemplate/pkg/lg"
	"gotemplate/pkg/mtr"
	"gotemplate/pkg/trc"
)

type Infrastructure struct {
	Logger                lg.Logger
	PrometheusErrorLogger lg.PrometheusErrorLogger
	Tracing               *trc.TracerProvider
	Metric                *mtr.MeterProvider
}

func Setup(configs *cfg.Config) (*Infrastructure, error) {
	infra := &Infrastructure{
		Logger: lg.New(&lg.Config{
			Level: configs.Log.Level,
			Mode:  configs.Log.Mode,
		}),
	}

	ctx := context.Background()
	if configs.Tracing.Kind != "" {
		tracing, err := trc.Setup(ctx, &trc.Config{
			Kind:        trc.Kind(configs.Tracing.Kind),
			ServiceName: configs.Tracing.ServiceName,
			Exporter: trc.ExporterConfig{
				Kind:     trc.ExporterKind(configs.Tracing.Exporter.Kind),
				Endpoint: configs.Tracing.Exporter.Endpoint,
			},
			Provider: trc.ProviderConfig{
				SamplerKind: trc.SamplerKind(configs.Tracing.Provider.SamplerKind),
				SampleRate:  configs.Tracing.Provider.SampleRate,
			},
		})
		if err != nil {
			infra.Logger.Warn("failed to setup otel tracing, continuing without tracing", lg.Err(err))
		} else {
			infra.Tracing = tracing
			infra.Logger.Info("tracing ready", lg.Str("kind", configs.Tracing.Kind))
		}
	}

	if configs.Metric.Kind != "" {
		metric, err := mtr.Setup(ctx, &mtr.Config{
			ServiceName: configs.Metric.ServiceName,
			Kind:        mtr.Kind(configs.Metric.Kind),
			Readers:     configs.Metric.Readers.ToMtrReaderConfigs(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to setup metric: %w", err)
		}
		infra.Metric = metric
		infra.Logger.Info("metric ready", lg.Str("kind", configs.Metric.Kind))
	}

	if infra.Metric != nil && infra.Metric.PromClient != nil {
		infra.PrometheusErrorLogger = lg.NewPrometheusErrorLogger(&lg.Config{
			Level: configs.Log.Level,
			Mode:  configs.Log.Mode,
		}, infra.Logger)
	}

	return infra, nil
}

func (i *Infrastructure) Close(ctx context.Context) error {
	if i.Tracing != nil {
		if err := i.Tracing.StopFn(ctx); err != nil {
			return fmt.Errorf("failed to stop tracing: %w", err)
		}
		i.Logger.Info("tracing shutdown complete")
	}

	if i.Metric != nil {
		if err := i.Metric.StopFn(ctx); err != nil {
			return fmt.Errorf("failed to stop metric: %w", err)
		}
		i.Logger.Info("metric shutdown complete")
	}
	return nil
}
