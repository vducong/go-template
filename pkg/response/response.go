package response

import (
	"net/http"
	"promotion/pkg/failure"

	"github.com/gin-gonic/gin"
)

const DefaultResponseSuccess = "success"

type HTTPResponse struct {
	Status failure.ErrCode `json:"status"`
	Data   any             `json:"data"`
}

func ErrBinding(ctx *gin.Context, err *failure.BindJSONErr) {
	_ = ctx.Error(err)
}

func ErrApp(ctx *gin.Context, err *failure.AppErr) {
	_ = ctx.Error(err)
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, HTTPResponse{
		Status: http.StatusOK,
		Data:   data,
	})
}
