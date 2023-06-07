package app

import (
	"promotion/configs"
	"promotion/internal/controller"
	"promotion/internal/middleware"
	"promotion/internal/router"
	"promotion/internal/server"
)

func Run(cfg *configs.Config) {
	infra := initInfra(cfg)
	engine := initServerDeps(cfg, infra)
	server.New(&server.CreateServerDTO{
		Cfg:     cfg,
		Log:     infra.Log,
		Handler: engine.Handler,
	})
}

func initServerDeps(cfg *configs.Config, infra *infrastructure) *router.Engine {
	authMiddlewares := middleware.New(cfg, infra.db.Firebase)
	modules := initModules(infra)
	controllers := initControllers(infra, modules)
	return router.NewEngine(&router.CreateEngineDTO{
		Cfg: cfg, Log: infra.Log, Controllers: controllers, AuthMiddlewares: authMiddlewares,
	})
}

func initControllers(infra *infrastructure, modules *Modules) *controller.Controllers {
	return &controller.Controllers{
		HealthCheck:  controller.NewHealthCheckController(),
		ReusableCode: controller.NewReusableCodeController(infra.Log, modules.ReusableCode),
	}
}
