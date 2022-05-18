package middleware

import (
	"net/http"
	"promotion/internal/dtos"
	"promotion/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	JWTAuthTokenPrefix = "Bearer"
	JWTAuthHeader      = "authorization"
)

type JWTAuthMiddleware struct {
	service *services.JWTAuthService
}

func NewJWTAuthMiddleware(
	service *services.JWTAuthService,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		service: service,
	}
}

func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(JWTAuthHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.HTTPResponse{
				Data: "missing auth header",
			})
			return
		}

		idToken := m.getIDToken(authHeader)
		if idToken == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.HTTPResponse{
				Data: "missing jwt",
			})
			return
		}

		token, err := m.service.VerifyIDToken(idToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.HTTPResponse{
				Data: "malformed jwt",
			})
			return
		}

		ctx.Set("UUID", token.UID)
		ctx.Next()
	}
}

func (m *JWTAuthMiddleware) getIDToken(header string) string {
	return strings.TrimSpace(strings.Replace(header, JWTAuthTokenPrefix, "", 1))
}
