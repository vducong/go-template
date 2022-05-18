package routers

import "promotion/internal/controller"

type Routers struct {
	healthCheck  *HealthCheckRouter
	reusableCode *ReusableCodeRouter
}

func (e *Engine) registerRouters(controllers *controller.Controllers) {
	rootHandler := e.Handler.Group("/promotion")
	e.routers = &Routers{
		healthCheck:  NewHealthCheckRouter(rootHandler, controllers.HealthCheck),
		reusableCode: NewReusableCodeRouter(rootHandler, controllers.ReusableCode),
	}
}
