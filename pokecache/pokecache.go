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
	go cache.reapLoop(interval)
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

func (C *Cache) Get(key string) ([]byte, bool){
	C.mu.Lock()
	defer C.mu.Unlock()

	entry, exists := C.data[key]
	if !exists {
		var dat []byte
		return dat, false
	}
	return entry.val, true
}

func (C *Cache) reapLoop(interval time.Duration) {
    for {
        time.Sleep(interval) 
        C.mu.Lock()
        for key, entry := range C.data {
            if time.Since(entry.createdAt) > 5*time.Second {
                delete(C.data, key)
            }
        }
        C.mu.Unlock()
    }
}
