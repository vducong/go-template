package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"promotion/pkg/logger"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				handleRecoverErr(ctx, log, err)
			}
		}()
		ctx.Next()
	}
}

func handleRecoverErr(ctx *gin.Context, log *logger.Logger, err interface{}) {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if netErr, ok := err.(*net.OpError); ok {
		if osErr, ok := netErr.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(osErr.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(osErr.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
	if brokenPipe {
		log.Error(
			ctx.Request.URL.Path,
			zap.Any("error", err),
			zap.String("request", string(httpRequest)),
		)
		// If the connection is dead, we can't write a status to it.
		ctx.Error(err.(error)) // nolint: errcheck
		ctx.Abort()
		return
	}

	log.Error(
		"[Recovery from panic]",
		zap.Time("time", time.Now()),
		zap.Any("error", err),
		zap.String("request", string(httpRequest)),
		zap.String("stack", string(debug.Stack())),
	)

	ctx.AbortWithStatus(http.StatusInternalServerError)
}
