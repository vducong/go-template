package middleware

import (
	"gotemplate/internal/cfg"
	"net/http"
	"strings"

	"github.com/go-chi/cors"
)

// CORS returns browser CORS middleware scoped to the public API (/api/v1).
// It must be mounted at the top-level router so preflight OPTIONS is handled;
// go-chi/cors does not work correctly inside Route groups alone.
// /api/v1/internal, /health, and /metrics are left unchanged.
func CORS(c cfg.CORSConfig) func(http.Handler) http.Handler {
	handler := cors.Handler(cors.Options{
		AllowedOrigins:   c.AllowedOrigins,
		AllowedMethods:   c.AllowedMethods,
		AllowedHeaders:   c.AllowedHeaders,
		ExposedHeaders:   c.ExposedHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
	})

	return func(next http.Handler) http.Handler {
		corsNext := handler(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicAPIV1(r.URL.Path) {
				corsNext.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isPublicAPIV1(path string) bool {
	if path == "/api/v1" || strings.HasPrefix(path, "/api/v1/") {
		return !isInternalAPIV1(path)
	}
	return false
}

func isInternalAPIV1(path string) bool {
	return path == "/api/v1/internal" || strings.HasPrefix(path, "/api/v1/internal/")
}
