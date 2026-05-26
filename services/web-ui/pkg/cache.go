package cache

import (
	"sync"
	"time"
)

type item[T any] struct {
	value     T
	expiresAt time.Time
}

type MemCache[T any] struct {
	mu    sync.RWMutex
	items map[string]item[T]
}

func NewMemCache[T any](cleanupInterval time.Duration) *MemCache[T] {
	c := &MemCache[T]{
		items: make(map[string]item[T]),
	}

	go c.startCleaner(cleanupInterval)
	return c
}

func (c *MemCache[T]) Set(key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item[T]{
		value: value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zero T

	item, exists := c.items[key]
	if !exists {
		return zero, false
	}

	if time.Now().After(item.expiresAt) {
		return zero, false
	}

	return item.value, true
}

func (c *MemCache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemCache[T]) startCleaner(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
