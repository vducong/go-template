package cache

import "fmt"

func GetCacheKey(template CacheKeyTemplate, params ...string) string {
	return fmt.Sprintf(string(template), params)
}
