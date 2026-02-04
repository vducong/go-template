package handler

type Config struct{}

type Handlers struct {
	HealthHandler *HealthHandler
}

func Setup(configs *Config) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
	}
}
