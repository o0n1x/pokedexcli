package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	var cache = Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.reapLoop()

	return &cache

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cacheval, ok := c.cache[key]
	if !ok {
		return nil, ok
	} else {
		fmt.Println("Cache used!")
		return cacheval.val, ok
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.cache {
			if time.Since(val.createdAt) > c.interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
