package httpif

import (
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/httpif/handler"
	"gotemplate/internal/interface/httpif/middleware"
	"gotemplate/pkg/httpsvr"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(
	configs *cfg.Config,
	infrastructure *infra.Infrastructure,
	handlers *handler.Handlers,
	authenticator *middleware.Authenticator,
) *httpsvr.Server {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	if infrastructure.Tracing != nil {
		r.Use(middleware.Tracing(configs.Tracing.ServiceName))
	}
	r.Use(middleware.RequestLogger(infrastructure.Logger))
	r.Use(middleware.Recovery(infrastructure.Logger, infrastructure.ResponseWriter))
	if configs.HTTP.CORS.Enabled {
		r.Use(middleware.CORS(configs.HTTP.CORS))
	}

	if infrastructure.Metric != nil && infrastructure.Metric.PromClient != nil {
		opts := promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}
		if infrastructure.PrometheusErrorLogger != nil {
			opts.ErrorLog = infrastructure.PrometheusErrorLogger
		}
		r.Handle("/metrics", promhttp.HandlerFor(infrastructure.Metric.PromClient.Registry, opts))
	}

	r.Get("/health", handlers.Health.Check)

	r.Route("/api/v1/internal", func(api chi.Router) {
		api.Use(authenticator.RequireAPIKey())
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(authenticator.RequireJwt())
	})

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
