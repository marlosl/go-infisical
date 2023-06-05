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
		Shards: 1024,
		LifeWindow: 1 * time.Hour,
		CleanWindow: 5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize: 500,
		Verbose: false,
		HardMaxCacheSize: 256,
		OnRemove: nil,
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
