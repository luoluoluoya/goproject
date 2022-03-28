package cache

import "sync"

type cache struct {
	mu         sync.RWMutex
	lru        *Cache
	cacheBytes int64
}

func (c *cache) Set(key string, value ByteView) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = NewCache(c.cacheBytes, nil)
	}
	return c.lru.Set(key, value)
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
