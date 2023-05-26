package controller

import (
	"promotion/pkg/failure"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	HealthCheck  *HealthCheckController
	ReusableCode *ReusableCodeController
}

func BindJSON[B interface{}](ctx *gin.Context) (bindedBody *B, err *failure.BindJSONErr) {
	var body B
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return nil, &failure.BindJSONErr{
			OriginalErr: err,
			Model:       reflect.TypeOf(body),
		}
	}
	return &body, nil
}
