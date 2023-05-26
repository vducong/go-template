package middleware

import (
	"net/http"
	"promotion/pkg/failure"
	"promotion/pkg/logger"
	"promotion/pkg/response"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		for _, err := range ctx.Errors {
			handleContextErr(ctx, log, err)
		}
	}
}

func handleContextErr(ctx *gin.Context, log *logger.Logger, err *gin.Error) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)

	if appError, ok := err.Err.(*failure.BindJSONErr); ok {
		if appError.OriginalErr != nil {
			log.Error(appError.OriginalErr)
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.HTTPResponse{
			Status: appError.Code,
			Data:   appError.Error(),
		})
		return
	}

	if appError, ok := err.Err.(*failure.AppErr); ok {
		if appError.OriginalErr != nil {
			log.Error(appError.OriginalErr)
		}
		ctx.AbortWithStatusJSON(appError.HTTPCode(), response.HTTPResponse{
			Status: appError.Code,
			Data:   appError.Error(),
		})
		return
	}

	log.Error(err.Err)
	ctx.JSON(http.StatusBadRequest, response.HTTPResponse{
		Status: http.StatusBadRequest,
		Data:   err.Error(),
	})
}
