package main

import (
	"gotemplate/internal/app"
	cfg "gotemplate/internal/config"
)

func main() {
	configs, err := cfg.Load("")
	if err != nil {
		panic("load config: " + err.Error())
	}

	a, err := app.Init(configs)
	if err != nil {
		panic("init app: " + err.Error())
	}

	a.Run()
}
