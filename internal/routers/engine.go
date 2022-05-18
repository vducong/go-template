package routers

import (
	"promotion/configs"
	"promotion/internal/controller"
	"promotion/internal/middleware"
	"promotion/pkg/logger"

	"github.com/gin-gonic/gin"
)

type CreateEngineDTO struct {
	Cfg         *configs.Config
	Log         *logger.Logger
	Controllers *controller.Controllers
}

type Engine struct {
	log     *logger.Logger
	Handler *gin.Engine
	routers *Routers
}

func NewEngine(dto *CreateEngineDTO) *Engine {
	if dto.Cfg.Server.Env == configs.ServerEnvProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	gin.ForceConsoleColor()

	engine := &Engine{
		log:     dto.Log,
		Handler: gin.New(),
	}

	gin.DebugPrintRouteFunc = logger.DebugOutputLogger(dto.Log)

	engine.attachMiddleware()
	engine.registerRouters(dto.Controllers)
	return engine
}

func (e *Engine) attachMiddleware() {
	e.Handler.Use(middleware.LoggerMiddleware(e.log))
	e.Handler.Use(middleware.RecoveryMiddleware(e.log))
}
