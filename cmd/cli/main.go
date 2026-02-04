package main

import (
	"flag"
	"fmt"
	"gotemplate/internal/app"
	cfg "gotemplate/internal/config"
	"os"
)

func main() {
	runCmd := flag.Bool("run", true, "run the application")
	configFileCmd := flag.String("config-file", "", "config file path")

	flag.Parse()

	if runCmd != nil && *runCmd {
		fmt.Println("running the application")
		if err := runApp(configFileCmd); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run app: %v\n", err)
			os.Exit(1)
		}
		return
	}
}

func runApp(configFilePath *string) error {
	path := ""
	if configFilePath != nil && *configFilePath != "" {
		path = *configFilePath
	}

	configs, err := cfg.Load(path)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	a, err := app.Init(configs)
	if err != nil {
		return fmt.Errorf("init app: %w", err)
	}

	a.Run()
	return nil
}
