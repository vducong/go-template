package cfld

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONConfigLoader struct{}

func (l *JSONConfigLoader) Load(path string, structure any) error {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read json config from file: %w", err)
	}

	return json.Unmarshal(jsonData, structure)
}
