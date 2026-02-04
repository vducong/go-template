package http

import (
	cfg "gotemplate/internal/config"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/http/handler"
	"gotemplate/internal/interface/http/middleware"
	"gotemplate/pkg/httpsvr"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(configs *cfg.Config, infrastructure *infra.Infrastructure, handlers *handler.Handlers) *httpsvr.Server {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	if infrastructure.Tracing != nil {
		r.Use(middleware.Tracing(configs.Tracing.ServiceName))
	}
	r.Use(middleware.RequestLogger(infrastructure.Logger))

	if infrastructure.Metric != nil && infrastructure.Metric.PromClient != nil {
		opts := promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}
		if infrastructure.PrometheusErrorLogger != nil {
			opts.ErrorLog = infrastructure.PrometheusErrorLogger
		}
		r.Handle("/metrics", promhttp.HandlerFor(infrastructure.Metric.PromClient.Registry, opts))
	}

	r.Get("/health", handlers.HealthHandler.Check)

	return httpsvr.New(
		httpsvr.WithConfig(&httpsvr.Config{
			Port:              configs.HTTP.Port,
			ReadTimeout:       configs.HTTP.ReadTimeout,
			ReadHeaderTimeout: configs.HTTP.ReadHeaderTimeout,
			WriteTimeout:      configs.HTTP.WriteTimeout,
			IdleTimeout:       configs.HTTP.IdleTimeout,
		}),
		httpsvr.WithHandler(r),
	)
}
