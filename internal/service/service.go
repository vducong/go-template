package service

import (
	"context"
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	"gotemplate/internal/repository"
)

type Services struct {
	Config Config
}

func Setup(
	configs *cfg.Config,
	infrastructure *infra.Infrastructure,
	repos *repository.Repos,
) (*Services, error) {
	configService, err := NewConfigService(configs)
	if err != nil {
		return nil, err
	}

	return &Services{
		Config: configService,
	}, nil
}

func (s *Services) Stop(ctx context.Context) error {
	if s == nil {
		return nil
	}

	return nil
}
