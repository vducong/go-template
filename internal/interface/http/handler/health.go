package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
