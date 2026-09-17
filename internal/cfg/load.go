package cfg

import (
	"fmt"
	"gotemplate/pkg/cfld"
)

func Load(path string) (*Config, error) {
	fmt.Printf("loading config from: %s\n", path)

	configs := &Config{}
	loader := cfld.New(cfld.LoaderTypeCleanenv)
	if err := loader.Load(path, configs); err != nil {
		return nil, fmt.Errorf("load config from %s: %w", path, err)
	}
	fmt.Printf("config loaded from %s: app=%s version=%s env=%s http_port=%s grpc_port=%s\n",
		path, configs.App.Name, configs.App.Version, configs.App.Env, configs.HTTP.Port, configs.GRPC.Port)

	return configs, nil
}
