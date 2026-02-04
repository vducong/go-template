package mtr

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type PromClient struct {
	Registry *prometheus.Registry
}

func setupProm() *PromClient {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
		),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			ReportErrors: true,
		}),
		collectors.NewBuildInfoCollector(),
	)

	return &PromClient{
		Registry: reg,
	}
}
