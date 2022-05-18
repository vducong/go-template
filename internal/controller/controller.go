package controller

import (
	"promotion/internal/services"
	"promotion/pkg/failure"
	"promotion/pkg/logger"
	"reflect"

	"github.com/gin-gonic/gin"
)

type CreateControlllersDTO struct {
	Log      *logger.Logger
	Services *services.Services
}

type Controllers struct {
	HealthCheck  *HealthCheckController
	ReusableCode *ReusableCodeController
}

func New(dto *CreateControlllersDTO) *Controllers {
	return &Controllers{
		HealthCheck:  NewHealthCheckController(),
		ReusableCode: NewReusableCodeController(dto.Log, dto.Services.ReusableCode),
	}
}

func BindJSON[B interface{}](ctx *gin.Context) (bindedBody *B, err string) {
	var body B
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return nil, failure.BindJSONErr{
			Model: reflect.TypeOf(body),
			Err:   err,
		}.String()
	}
	return &body, ""
}
