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
		return nil, fmt.Errorf("failed to load server config: %w", err)
	}
	fmt.Printf("server config loaded: %+v\n", configs)

	return configs, nil
}
