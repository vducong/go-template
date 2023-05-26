package controller

import (
	"promotion/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckController struct{}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) HealthCheck(startedAt time.Time) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		uptime := time.Since(startedAt)
		response.Success(ctx, gin.H{
			"started_at": startedAt.String(),
			"uptime":     uptime.String(),
			"ip_address": ctx.ClientIP(),
		})
	}
}
