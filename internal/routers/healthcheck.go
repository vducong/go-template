package routers

import (
	"promotion/internal/controller"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckRouter struct {
	handler    *gin.RouterGroup
	controller *controller.HealthCheckController
}

func NewHealthCheckRouter(
	handler *gin.RouterGroup,
	ctrler *controller.HealthCheckController,
) *HealthCheckRouter {
	router := &HealthCheckRouter{
		handler:    handler,
		controller: ctrler,
	}
	router.Handle()
	return router
}

func (r *HealthCheckRouter) Handle() {
	r.handler.GET("/healthcheck", r.controller.HealthCheck(time.Now()))
}
