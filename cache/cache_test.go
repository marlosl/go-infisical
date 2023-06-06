package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache, err := NewCache()
	assert.NoError(t, err)

	// Test Update and Read
	err = cache.Update("key1", "value1")
	assert.NoError(t, err)

	value, err := cache.Read("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)

	// Test Delete
	cache.Delete("key1")

	value, err = cache.Read("key1")
	assert.Error(t, err)
	assert.Equal(t, "", value)

	// Test expiration
	err = cache.Update("key2", "value2")
	assert.NoError(t, err)

	time.Sleep(2 * time.Hour)

	value, err = cache.Read("key2")
	assert.Error(t, err)
	assert.Equal(t, "", value)
}
