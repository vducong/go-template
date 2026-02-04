package service

import cfg "gotemplate/internal/config"

type Services struct {
	Config Config
}

func Setup(configs *cfg.Config) (*Services, error) {
	configService, err := NewConfigService(configs)
	if err != nil {
		return nil, err
	}

	return &Services{
		Config: configService,
	}, nil
}
