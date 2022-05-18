package app

import (
	"promotion/configs"
	"promotion/internal/controller"
	"promotion/internal/middleware"
	"promotion/internal/repos"
	"promotion/internal/routers"
	"promotion/internal/server"
	"promotion/internal/services"
	"promotion/pkg/databases"
	"promotion/pkg/logger"
	"promotion/pkg/pubsub"
)

func Run(cfg *configs.Config) {
	log := logger.New(cfg)

	db, err := databases.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize DB")
	}
	defer db.Close()

	if _, err := pubsub.NewPubSub(cfg, log); err != nil {
		log.Fatalf("Failed to initialize PubSub")
	}

	repo := repos.New(&db)

	service := services.New(&services.CreateServicesDTO{
		Cfg: cfg, Log: log, DB: &db, Repos: repo,
	})

	controllers := controller.New(&controller.CreateControlllersDTO{
		Log: log, Services: service,
	})

	middleware.New(&middleware.CreateMiddlewaresDTO{
		Cfg: cfg, JWTAuthService: service.JWTAuth,
	})

	engine := routers.NewEngine(&routers.CreateEngineDTO{
		Cfg: cfg, Log: log, Controllers: controllers,
	})

	server.New(&server.CreateServerDTO{
		Cfg: cfg, Log: log, Handler: engine.Handler,
	})
}
