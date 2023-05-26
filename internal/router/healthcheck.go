package router

import (
	"promotion/internal/controller"
	"time"

	"github.com/gin-gonic/gin"
)

type healthCheckRouter struct {
	group      *gin.RouterGroup
	controller *controller.HealthCheckController
}

func initHealthCheckRouter(
	group *gin.RouterGroup,
	c *controller.HealthCheckController,
) {
	router := newHealthCheckRouter(group, c)
	router.handle()
}

func newHealthCheckRouter(
	group *gin.RouterGroup,
	c *controller.HealthCheckController,
) *healthCheckRouter {
	return &healthCheckRouter{group, c}
}

func (r healthCheckRouter) handle() {
	root := r.group.Group("/health")
	root.GET("", r.controller.HealthCheck(time.Now()))
}
