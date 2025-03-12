package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Cachemap map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
	stopChan chan struct{}
}

func NewCache(duration time.Duration) *Cache {
	newcache := &Cache{
		Cachemap: make(map[string]cacheEntry),
		interval: duration,
		stopChan: make(chan struct{}),
	}
	go newcache.reapLoop(duration)
	return newcache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Cachemap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.Cachemap[key]
	if !exists {
		return nil, false
	}
	return value.val, true
}

func (c *Cache) reapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case ticktime := <-ticker.C:
			c.mu.Lock()
			keystodelete := []string{}
			for key, value := range c.Cachemap {
				time_passed := ticktime.Sub(value.createdAt)
				if time_passed > duration {
					keystodelete = append(keystodelete, key)
				}
			}
			c.mu.Unlock()

			c.mu.Lock()
			for _, key := range keystodelete {
				delete(c.Cachemap, key)
			}
			c.mu.Unlock()
		// currently not using stop feature
		case <-c.stopChan:
			return
		}
	}
}

// currently not in use
func (c *Cache) Close() {
	close(c.stopChan)
}
