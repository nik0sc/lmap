// Package lru contains a simple implementation of a LRU cache.
// It is inspired by the Cache type at https://github.com/hashicorp/golang-lru.
package lru

import (
	"sync"

	"github.com/nik0sc/lmap"
)

const (
	DefaultCacheMax = 100
)

// Cache is a LRU cache. It is safe for concurrent use.
// Only one type of key may be stored.
type Cache struct {
	l   *lmap.LinkedMap
	mu  sync.Mutex
	max int
}

// New creates a new Cache ready for use.
// If max > 0, the cache will only be allowed to contain max number of entries.
// Otherwise a default maximum number will be used.
func New(max int) Cache {
	if max <= 0 {
		max = DefaultCacheMax
	}

	return Cache{
		l:   lmap.New(),
		max: max,
	}
}

// Add adds a key-value pair to the cache.
// If a cache entry is evicted as a result, Add returns true.
// If the type of key is inconsistent with previous cache additions,
// Add panics.
func (c *Cache) Add(key, value interface{}) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.l.Get(key, false)

	if !ok && c.max > 0 && c.l.Len() >= c.max {
		evicted = true
		c.l.Head(true)
	}

	c.l.Set(key, value, true)
	return
}

// Get reads a value from the cache.
// If the key was not found, ok will be false.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.l.Get(key, true)
}

// Peek is like Get, but it will not affect the recency of the key.
func (c *Cache) Peek(key interface{}) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.l.Get(key, false)
}

// Len returns the number of elements in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	l := c.l.Len()
	c.mu.Unlock()
	return l
}

// Trim removes least recently used elements from the cache
// so that its size is at most max elements.
func (c *Cache) Trim(max int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for c.l.Len() > max {
		c.l.Head(true)
	}
}

// Keys returns the keys in the cache, ordered by increasing recency.
func (c *Cache) Keys() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	ks := make([]interface{}, 0, c.l.Len())
	c.l.Iter(func(k, _ interface{}) bool {
		ks = append(ks, k)
		return true
	})

	return ks
}
