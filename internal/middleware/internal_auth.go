package middleware

import (
	"errors"
	"net/http"
	"promotion/configs"
	"promotion/pkg/response"

	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
)

const internalAuthHeader = "x-api-key"

type InternalAuthMiddleware struct {
	cfg *configs.Config
}

func NewInternalAuthMiddleware(cfg *configs.Config) *InternalAuthMiddleware {
	return &InternalAuthMiddleware{cfg: cfg}
}

func (m *InternalAuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if errHeader := verifyAuthHeader(ctx, m.cfg); errHeader == nil {
			extractTracingParentSpanIntoCtx(ctx)
			ctx.Next()
			return
		}

		errQuery := verifyAuthTokenFromQuery(ctx, m.cfg)
		if errQuery == nil {
			extractTracingParentSpanIntoCtx(ctx)
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPResponse{
			Status: http.StatusUnauthorized,
			Data:   errQuery.Error(),
		})
		return
	}
}

func verifyAuthHeader(ctx *gin.Context, cfg *configs.Config) error {
	authHeader := ctx.GetHeader(internalAuthHeader)
	if authHeader == "" {
		return errors.New("Missing auth header")
	}

	if authHeader != cfg.APIKey.PromotionAPIKey {
		return errors.New("Invalid auth")
	}
	return nil
}

func verifyAuthTokenFromQuery(ctx *gin.Context, cfg *configs.Config) error {
	// Mostly for the use of PubSub
	authToken, ok := ctx.GetQuery(internalAuthHeader)
	if !ok {
		return errors.New("Missing auth header")
	}
	if authToken != cfg.APIKey.PromotionAPIKey {
		return errors.New("Invalid auth")
	}
	return nil
}

func extractTracingParentSpanIntoCtx(ctx *gin.Context) {
	propagator := gcppropagator.CloudTraceFormatPropagator{}
	reqCtx := propagator.Extract(
		ctx.Request.Context(), propagation.HeaderCarrier(ctx.Request.Header),
	)
	ctx.Request = ctx.Request.Clone(reqCtx)
}
