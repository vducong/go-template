package httprespwrit

import (
	"context"
	"net/http"
)

type Logger interface {
	Error(err error)
	CtxError(ctx context.Context, err error)
}

type Writer interface {
	WriteJSON(w http.ResponseWriter, req *http.Request, res *JSONResponse)
	WriteError(w http.ResponseWriter, req *http.Request, res *ErrorResponse)
}

func NewWriter(l Logger) Writer {
	return newRenderWriter(l)
}
