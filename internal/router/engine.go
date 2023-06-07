package router

import (
	"promotion/configs"
	"promotion/internal/controller"
	"promotion/internal/middleware"
	"promotion/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type CreateEngineDTO struct {
	Cfg             *configs.Config
	Log             *logger.Logger
	Controllers     *controller.Controllers
	AuthMiddlewares *middleware.AuthMiddlewares
}

type Engine struct {
	log     *logger.Logger
	Handler *gin.Engine
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

	engine.attachMiddleware(dto.Cfg)
	engine.registerRoutes(dto.Controllers, dto.AuthMiddlewares)
	return engine
}

func (e *Engine) attachMiddleware(cfg *configs.Config) {
	e.Handler.Use(middleware.ErrorHandler(e.log))
	e.Handler.Use(middleware.RecoveryMiddleware(e.log))
	e.Handler.Use(otelgin.Middleware(cfg.Server.Name))

	if cfg.Server.Env == configs.ServerEnvLocalhost {
		e.Handler.Use(middleware.LoggerMiddleware(e.log))
	}
}

func (e *Engine) registerRoutes(
	controllers *controller.Controllers,
	authMiddlewares *middleware.AuthMiddlewares,
) {
	root := e.Handler.Group("promotion")
	initHealthCheckRouter(root, controllers.HealthCheck)
	initReusableCodeRouter(root, controllers.ReusableCode, authMiddlewares.InternalAuth)
}
