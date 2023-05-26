package app

import (
	"promotion/configs"
	"promotion/pkg/databases"
	"promotion/pkg/logger"
	"promotion/pkg/pubsub"
	"promotion/pkg/tracing"
)

type infrastructure struct {
	log    *logger.Logger
	db     databases.Databases
	ps     *pubsub.PubSub
	tracer *tracing.Tracer
}

func initInfra(cfg *configs.Config) *infrastructure {
	log := logger.New(cfg)

	db, err := databases.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize DB")
	}

	ps, err := pubsub.NewPubSub(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize PubSub")
	}
	log.Info("PubSub initialized successfully")

	tracer, err := tracing.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to init trace exporter: %v", err)
	}
	log.Info("Tracer initialized successfully")

	return &infrastructure{log, db, ps, tracer}
}
