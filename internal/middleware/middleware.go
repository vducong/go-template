package middleware

import (
	"promotion/configs"
	"promotion/pkg/databases"
)

type AuthMiddlewares struct {
	InternalAuth *InternalAuthMiddleware
	JWTAuth      *JWTAuthMiddleware
}

func New(cfg *configs.Config, fbAuth databases.FirebaseAuth) *AuthMiddlewares {
	return &AuthMiddlewares{
		InternalAuth: NewInternalAuthMiddleware(cfg),
		JWTAuth:      NewJWTAuthMiddleware(fbAuth),
	}
}
