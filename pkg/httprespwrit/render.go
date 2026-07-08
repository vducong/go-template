package httprespwrit

import (
	"net/http"

	"github.com/go-chi/render"
)

type RenderWriter struct {
	logger Logger
}

func newRenderWriter(l Logger) Writer {
	return &RenderWriter{
		logger: l,
	}
}

func (rw *RenderWriter) WriteJSON(w http.ResponseWriter, req *http.Request, res *JSONResponse) {
	if err := render.Render(w, req, res); err != nil {
		rw.logger.CtxError(req.Context(), err)
	}
}

func (rw *RenderWriter) WriteError(w http.ResponseWriter, req *http.Request, res *ErrorResponse) {
	if res.Err != nil {
		rw.logger.CtxError(req.Context(), res.Err)
	}
	if err := render.Render(w, req, res); err != nil {
		rw.logger.CtxError(req.Context(), err)
	}
}
