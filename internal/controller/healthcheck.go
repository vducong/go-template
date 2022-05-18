package controller

import (
	"net/http"
	"promotion/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	Log *logger.Logger
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) HealthCheck(startedAt time.Time) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		uptime := time.Since(startedAt)
		ctx.JSON(http.StatusOK, gin.H{
			"started_at": startedAt.String(),
			"uptime":     uptime.String(),
			"ip_address": ctx.ClientIP(),
		})
	}
}
