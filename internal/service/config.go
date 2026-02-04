package service

import (
	cfg "gotemplate/internal/config"
)

type Config interface {
	Get() *cfg.Config
}

type ConfigService struct {
	configs *cfg.Config
}

func NewConfigService(configs *cfg.Config) (Config, error) {
	return &ConfigService{
		configs: configs,
	}, nil
}

func (s *ConfigService) Get() *cfg.Config {
	return s.configs
}
