package httpif

import (
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/httpif/handler"
	"gotemplate/internal/interface/httpif/middleware"
	"gotemplate/pkg/httpsvr"
	"net/http"

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
	return httpsvr.New(
		httpsvr.WithConfig(&httpsvr.Config{
			Port:              configs.HTTP.Port,
			ReadTimeout:       configs.HTTP.ReadTimeout,
			ReadHeaderTimeout: configs.HTTP.ReadHeaderTimeout,
			WriteTimeout:      configs.HTTP.WriteTimeout,
			IdleTimeout:       configs.HTTP.IdleTimeout,
		}),
		httpsvr.WithHandler(newRouter(configs, infrastructure, handlers, authenticator)),
	)
}

func newRouter(
	configs *cfg.Config,
	infrastructure *infra.Infrastructure,
	handlers *handler.Handlers,
	authenticator *middleware.Authenticator,
) http.Handler {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
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

	// Tracing and request logging live on /api, not the router, so /metrics scrapes and /health probes skip them.
	// Recovery sits inside RequestLogger, so a panic is logged with the 500 it returns.
	r.Route("/api", func(api chi.Router) {
		if infrastructure.Tracing != nil {
			api.Use(middleware.Tracing(configs.Tracing.ServiceName))
		}
		api.Use(middleware.RequestLogger(infrastructure.Logger))
		api.Use(middleware.Recovery(infrastructure.Logger, infrastructure.ResponseWriter))

		api.Route("/v1/internal", func(internal chi.Router) {
			internal.Use(authenticator.RequireAPIKey())
		})

		api.Route("/v1", func(v1 chi.Router) {
			v1.Use(authenticator.RequireJwt())
		})
	})

	return r
}
