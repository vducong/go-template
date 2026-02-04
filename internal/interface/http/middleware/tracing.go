package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Tracing(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(
			next,
			serviceName,
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
		)
	}
}

// spanNameFormatter formats the span name using the HTTP method and chi route pattern.
// This provides more meaningful span names like "GET /users/{id}" instead of "GET /users/123".
func spanNameFormatter(_ string, r *http.Request) string {
	routePattern := chi.RouteContext(r.Context()).RoutePattern()
	if routePattern == "" {
		routePattern = r.URL.Path
	}
	return r.Method + " " + routePattern
}
