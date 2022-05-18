package controller

import (
	"net/http"
	"promotion/internal/dtos"
	"promotion/internal/services"
	"promotion/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ReusableCodeController struct {
	log     *logger.Logger
	service *services.ReusableCodeService
}

func NewReusableCodeController(
	log *logger.Logger, service *services.ReusableCodeService,
) *ReusableCodeController {
	return &ReusableCodeController{log, service}
}

func (c *ReusableCodeController) GetByCode(ctx *gin.Context) {
	body, errBinding := BindJSON[dtos.ReusableCodeGetByCodeReq](ctx)
	if errBinding != "" {
		c.log.Errorf("Binding JSON failed: %s", errBinding)
		ctx.JSON(http.StatusBadRequest, dtos.HTTPResponse{
			Data: errBinding,
		})
		return
	}

	rc, err := c.service.GetByCode(ctx, body.Code)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, dtos.HTTPResponse{
			Data: err,
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.HTTPResponse{
		Data: rc,
	})
}
