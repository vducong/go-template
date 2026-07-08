package apperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// Set once in Setup() at startup, read-only thereafter (no locking needed).
var (
	codePrefix    int
	messagePrefix string
	messageByKey  map[string]string
)

func Setup(prefix int, prefixKey string) error {
	if prefixKey == "" {
		return errors.New("message key prefix is required")
	}

	codePrefix = prefix
	messagePrefix = prefixKey

	catalog := make(map[string]string)
	if err := json.Unmarshal(messageCatalogJSON, &catalog); err != nil {
		return fmt.Errorf("parse embedded apperr messages: %w", err)
	}

	missing := make([]string, 0)
	for _, key := range orderedErrorKeys() {
		fullKey := messagePrefix + "." + key
		if _, ok := catalog[fullKey]; !ok {
			missing = append(missing, fullKey)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		return fmt.Errorf("apperr catalog missing keys for prefix %q: %v", messagePrefix, missing)
	}

	messageByKey = catalog
	return nil
}

func messageForKey(key string) (string, bool) {
	detail, ok := messageByKey[key]
	return detail, ok
}
