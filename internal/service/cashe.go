package service

import (
	"sync"
)

type Cache interface {
	Add(key, value string)
	Get(key string) (value string, ok bool)
	Len() int
}

type CacheMemory struct {
	mu    *sync.RWMutex
	items map[string]string
}

func NewCashe() Cache {
	return &CacheMemory{
		mu:    &sync.RWMutex{},
		items: make(map[string]string, 1000),
	}
}

func (c *CacheMemory) Add(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= 1000 {
		c.items = make(map[string]string, 1000)
	}

	c.items[key] = value
}

func (c *CacheMemory) Get(key string) (value string, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok = c.items[key]

	return
}

func (c *CacheMemory) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.items)
}
