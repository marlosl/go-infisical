package cache

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

type Cache struct {
	values *bigcache.BigCache
}

func NewCache() (*Cache, error) {
	bCache, err := bigcache.NewBigCache(bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 1 * time.Hour,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: false,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 256,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("new big cache: %w", err)
	}

	return &Cache{
		values: bCache,
	}, nil
}

func (c *Cache) Update(id, value string) error {
	return c.values.Set(id, []byte(value))
}

func (c *Cache) Read(id string) (string, error) {
	value, err := c.values.Get(id)
	if err != nil {
		return "", fmt.Errorf("get: %w", err)
	}

	return string(value), nil
}

func (c *Cache) Delete(id string) {
	c.values.Delete(id)
}
