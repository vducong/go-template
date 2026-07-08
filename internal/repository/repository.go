package repository

import (
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
)

type Repos struct {
	User UserStore
}

func Init(configs *cfg.Config, i *infra.Infrastructure) *Repos {
	return &Repos{
		User: NewUserStore(),
	}
}
