package main

import (
	"log"
	"promotion/configs"
	"promotion/internal/app"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	app.Run(cfg)
}
