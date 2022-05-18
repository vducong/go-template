package middleware

import (
	"promotion/internal/constants"
	"promotion/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()

		if len(ctx.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range ctx.Errors.Errors() {
				log.Error(e)
			}
		} else {
			end := time.Now()
			latency := end.Sub(start)
			fields := []zapcore.Field{
				zap.Int("status", ctx.Writer.Status()),
				zap.String("method", ctx.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", ctx.ClientIP()),
				zap.String("user-agent", ctx.Request.UserAgent()),
				zap.Duration("latency", latency),
				zap.String("time", end.Format(constants.DefaultTimeLayout)),
			}

			log.Desugar().Info(path, fields...)
		}
	}
}
