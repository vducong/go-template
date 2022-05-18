package middleware

import (
	"net/http"
	"promotion/configs"
	"promotion/internal/dtos"

	"github.com/gin-gonic/gin"
)

const internalAuthHeader = "x-api-key"

type InternalAuthMiddleware struct {
	cfg *configs.Config
}

func NewInternalAuthMiddleware(cfg *configs.Config) *InternalAuthMiddleware {
	return &InternalAuthMiddleware{
		cfg: cfg,
	}
}

func (m *InternalAuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(internalAuthHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.HTTPResponse{
				Data: "missing auth header",
			})
			return
		}

		if authHeader != m.cfg.APIKey.PromotionAPIKey {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.HTTPResponse{
				Data: "malformed auth header",
			})
			return
		}

		ctx.Next()
	}
}
