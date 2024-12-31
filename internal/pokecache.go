package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

type Cache struct{
	data map[string]cacheEntry
	mu sync.Mutex
}

func NewCache(interval time.Duration) *Cache{
	cache := &Cache{
		data: make(map[string]cacheEntry),
	}
	go func ()  {
		time.Sleep(interval)
	}()
	return cache;
}

func (C *Cache) Add(key string, val []byte) {
	C.mu.Lock()
	defer C.mu.Unlock()

	C.data[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	
}
