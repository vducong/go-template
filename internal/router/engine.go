package router

import (
	"promotion/configs"
	"promotion/internal/controller"
	"promotion/internal/middleware"
	"promotion/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Engine struct {
	log     *logger.Logger
	Handler *gin.Engine
}

func NewEngine(
	cfg *configs.Config,
	log *logger.Logger,
	controllers *controller.Controllers,
	authMiddlewares *middleware.AuthMiddlewares,
) *Engine {
	if cfg.Server.Env == configs.ServerEnvProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := &Engine{
		log:     log,
		Handler: gin.New(),
	}

	gin.ForceConsoleColor()
	gin.DebugPrintRouteFunc = logger.DebugOutputLogger(log)

	engine.attachMiddleware(cfg)
	engine.registerRoutes(controllers, authMiddlewares)
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
