package lru

import (
	"sync"

	"playground-1.16/dsa/lmap"
)

type Cache struct {
	l   *lmap.LinkedMap
	mu  *sync.Mutex
	max int
}

func New(max int) Cache {
	if max < 0 {
		max = 0
	}

	return Cache{
		l:   lmap.New(),
		mu: &sync.Mutex{},
		max: max,
	}
}

func (c Cache) Add(key, value interface{}) (evicted bool) {
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

func (c Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.l.Get(key, true)
}
