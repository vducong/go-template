package handler

import (
	"gotemplate/pkg/httprespwrit"
	"net/http"
)

type HealthHandler struct {
	ResponseWriter httprespwrit.Writer
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	h.ResponseWriter.WriteJSON(w, r, &httprespwrit.JSONResponse{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       "OK",
	})
}
