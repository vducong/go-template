package app

import (
	"promotion/configs"
	"promotion/pkg/cache"
	"promotion/pkg/databases"
	"promotion/pkg/logger"
	"promotion/pkg/pubsub"
	"promotion/pkg/tracing"
	"time"
)

type infrastructure struct {
	Log           *logger.Logger
	db            databases.Databases
	inMemoryCache *cache.InMemoryCache
	ps            *pubsub.PubSub
	tracer        *tracing.Tracer
}

func initInfra(cfg *configs.Config) *infrastructure {
	log := logger.New(cfg)

	db, err := databases.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize DB")
	}

	inMemoryCache, err := cache.NewInMemoryCache(&cache.CreateInMemoryCacheDTO{
		Log: log, TTL: time.Minute * 10,
	})
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}
	log.Info("In-memory cache initialized successfully")

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

	return &infrastructure{log, db, inMemoryCache, ps, tracer}
}
