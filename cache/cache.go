package cache

import (
	"sync"
	"time"
)

// Single Item type
type item[V any] struct {
	value      V
	expiration time.Time
}

// Cache type
type Cache[K comparable, V any] struct {
	items         map[K]item[V] // map of items
	mu            sync.Mutex
	cleanInterval time.Duration
}

const (
	NoExpiration time.Duration = -1
)

func (i item[V]) isExpired() bool {
	if i.expiration == (time.Time{}) {
		return false
	}
	return i.expiration.Before(time.Now())
}

// Create a new Cache with a background cleanup goroutine
func New[K comparable, V any](cleanInterval time.Duration) *Cache[K, V] {
	c := &Cache[K, V]{
		items: make(map[K]item[V]),
	}

	go func() {
		for range time.Tick(cleanInterval) {
			c.mu.Lock()
			for k, v := range c.items {
				if v.isExpired() {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

// Set a key value pair in the cache, ttl can be NoExpiration
func (c *Cache[K, V]) Set(key K, val V, ttl time.Duration) {
	var ex time.Time
	if ttl != NoExpiration {
		ex = time.Now().Add(ttl)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item[V]{
		value:      val,
		expiration: ex,
	}
}

// Get a value from the cache, returns value and a bool indicating if the key was valid
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return item.value, false // Zero value
	}
	if item.isExpired() {
		delete(c.items, key)
		return item.value, false
	}
	return item.value, true
}

// Delete a key from the cache
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Pop a key from the cache, returns value and a bool indicating if the key was valid
func (c *Cache[K, V]) Pop(key K) (V, bool) {
	val, has := c.Get(key)
	if !has {
		return val, has
	}
	c.Delete(key)
	return val, has
}
