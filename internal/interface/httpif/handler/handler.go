package handler

import (
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	intsvc "gotemplate/internal/service"
)

type Handlers struct {
	Health *HealthHandler
}

func Setup(
	configs *cfg.Config, infrastructure *infra.Infrastructure, services *intsvc.Services,
) *Handlers {
	return &Handlers{
		Health: &HealthHandler{
			ResponseWriter: infrastructure.ResponseWriter,
		},
	}
}
