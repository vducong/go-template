package routers

import (
	"promotion/internal/controller"

	"github.com/gin-gonic/gin"
)

type ReusableCodeRouter struct {
	handler    *gin.RouterGroup
	controller *controller.ReusableCodeController
}

func NewReusableCodeRouter(
	handler *gin.RouterGroup,
	ctrler *controller.ReusableCodeController,
) *ReusableCodeRouter {
	router := &ReusableCodeRouter{
		handler:    handler,
		controller: ctrler,
	}
	router.Handle()
	return router
}

func (r *ReusableCodeRouter) Handle() {
	r.handler.POST("/reusable-codes", r.controller.GetByCode)
}
