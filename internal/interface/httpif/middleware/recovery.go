package middleware

import (
	"gotemplate/internal/apperr"
	"gotemplate/pkg/httprespwrit"
	"gotemplate/pkg/lg"
	"net/http"
	"runtime/debug"
)

// Recovery catches panics, logs them, and returns a 500 JSON error.
// http.ErrAbortHandler is re-panicked so net/http can abort without treating it as a server fault.
func Recovery(logger lg.Logger, writer httprespwrit.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				if rec == http.ErrAbortHandler {
					panic(rec)
				}

				logger.Error("recovery middleware panicked",
					lg.Any("panic", rec),
					lg.Str("stack", string(debug.Stack())))

				writer.WriteError(w, r, &httprespwrit.ErrorResponse{
					Err: apperr.InternalServer,
				})
			}()
			next.ServeHTTP(w, r)
		})
	}
}
