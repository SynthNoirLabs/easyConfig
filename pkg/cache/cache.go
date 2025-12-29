package cache

import (
	"sync"
	"time"
)

// Cache is a thread-safe in-memory cache with TTL support.
type Cache struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
}

// CacheEntry represents an entry in the cache.
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// New creates a new Cache.
func New() *Cache {
	return &Cache{
		data: make(map[string]CacheEntry),
	}
}

// Get retrieves an item from the cache.
// It returns the value, a boolean indicating if the item was found,
// and a boolean indicating if the item is stale.
func (c *Cache) Get(key string) (interface{}, bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.data[key]
	if !found {
		return nil, false, false
	}

	stale := time.Now().After(entry.ExpiresAt)
	return entry.Data, true, stale
}

// Set adds an item to the cache with a specified TTL.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
