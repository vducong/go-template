package router

import (
	"promotion/internal/controller"
	"promotion/internal/middleware"

	"github.com/gin-gonic/gin"
)

type reusableCodeRouter struct {
	group        *gin.RouterGroup
	controller   *controller.ReusableCodeController
	internalAuth *middleware.InternalAuthMiddleware
}

func initReusableCodeRouter(
	group *gin.RouterGroup,
	c *controller.ReusableCodeController,
	internalAuth *middleware.InternalAuthMiddleware,
) {
	router := newReusableCodeRouter(group, c, internalAuth)
	router.handle()
}

func newReusableCodeRouter(
	group *gin.RouterGroup,
	c *controller.ReusableCodeController,
	internalAuth *middleware.InternalAuthMiddleware,
) *reusableCodeRouter {
	return &reusableCodeRouter{group, c, internalAuth}
}

func (r reusableCodeRouter) handle() {
	root := r.group.Group("/reusable-code")
	root.POST("", r.controller.GetByCode)
}
