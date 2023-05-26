package middleware

import (
	"net/http"
	"promotion/pkg/logger"
	timeutil "promotion/pkg/time_util"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)
		fields := []zapcore.Field{
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("ip", ctx.ClientIP()),
			zap.String("protocol", ctx.Request.Proto),
			zap.Duration("latency", latency),
			zap.Int("status", ctx.Writer.Status()),
			zap.String("time", end.Format(timeutil.DefaultTimeLayout)),
		}

		if ctx.Writer.Status() >= http.StatusMultipleChoices {
			log.Desugar().Error(path, fields...)
		} else {
			log.Desugar().Info(path, fields...)
		}
	}
}
