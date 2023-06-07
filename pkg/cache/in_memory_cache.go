package cache

import (
	"encoding/json"
	"promotion/pkg/failure"
	"promotion/pkg/logger"
	"time"

	"github.com/allegro/bigcache/v3"
)

type CreateInMemoryCacheDTO struct {
	Log *logger.Logger
	TTL time.Duration
}

type InMemoryCache struct {
	log   *logger.Logger
	cache *bigcache.BigCache
}

func NewInMemoryCache(dto *CreateInMemoryCacheDTO) (*InMemoryCache, error) {
	config := getConfig(dto.TTL)
	cache, err := bigcache.NewBigCache(*config)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &InMemoryCache{
		log:   dto.Log,
		cache: cache,
	}, nil
}

func (c *InMemoryCache) Close() {
	c.cache.Close()
}

func getConfig(ttl time.Duration) *bigcache.Config {
	return &bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: ttl,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 1 * time.Second,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed
		// because of its expiration time or no space left for the new entry, or because delete was called.
		// A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	}
}

func (c *InMemoryCache) Set(key string, value interface{}) error {
	cacheDataBytes, err := json.Marshal(value)
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	if err := c.cache.Set(key, cacheDataBytes); err != nil {
		return failure.ErrWithTrace(err)
	}
	return nil
}

func Get[DestinationType interface{}](c *InMemoryCache, key string) (*DestinationType, error) {
	valueInBytes, err := c.cache.Get(key)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	var value DestinationType
	if err := json.Unmarshal(valueInBytes, &value); err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	return &value, nil
}

func (c *InMemoryCache) Delete(key string) error {
	if err := c.cache.Delete(key); err != nil {
		return failure.ErrWithTrace(err)
	}
	return nil
}
