package cfld

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type CleanenvConfigLoader struct{}

func (l *CleanenvConfigLoader) Load(path string, structure any) error {
	if path == "" {
		path = os.Getenv(envPath)
	}
	if path == "" {
		path = defaultEnvPath
	}

	if path == "" {
		fmt.Println("reading config from environment variable only")
		if err := cleanenv.ReadEnv(structure); err != nil {
			return fmt.Errorf("read config from env: %w", err)
		}
		fmt.Println("config loaded from env")
		return nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "failed to read config from %s\n", path)
		fmt.Println("trying to read from env")
		if err := cleanenv.ReadEnv(structure); err != nil {
			return fmt.Errorf("read config from env: %w", err)
		}
		fmt.Println("config loaded")
		return nil
	}

	if err := cleanenv.ReadConfig(path, structure); err != nil {
		return fmt.Errorf("failed to read config from file: %w", err)
	}

	fmt.Printf("config loaded from file: %s\n", path)

	// override config from file with config from env
	if err := cleanenv.ReadEnv(structure); err != nil {
		fmt.Fprintf(os.Stderr, "could not load extra config from env: %v\n", err)
	}
	fmt.Println("extra config loaded from env")

	return nil
}
