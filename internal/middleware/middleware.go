package middleware

import (
	"promotion/configs"
	"promotion/internal/services"
)

type CreateMiddlewaresDTO struct {
	Cfg            *configs.Config
	JWTAuthService *services.JWTAuthService
}

type Middlewares struct {
	internalAuth *InternalAuthMiddleware
	jwtAuth      *JWTAuthMiddleware
}

func New(dto *CreateMiddlewaresDTO) *Middlewares {
	return &Middlewares{
		internalAuth: NewInternalAuthMiddleware(dto.Cfg),
		jwtAuth:      NewJWTAuthMiddleware(dto.JWTAuthService),
	}
}
