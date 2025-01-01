package pokecache

import (
	"fmt"
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
	go func() {
		cache.reapLoop()
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

func (C *Cache) Get(key string) ([]byte, bool){
	fmt.Println("41, poke.go")

	C.mu.Lock()
	defer C.mu.Unlock()
	fmt.Println("45, poke.go")

	entry, exists := C.data[key]

	fmt.Println("49, poke.go")

	if !exists {
		var dat []byte
		return dat, false
	}
	return entry.val, true
}

func (C *Cache) reapLoop() {
	for key, entry := range C.data{
		age := time.Since(entry.createdAt) 
		if age > 5 * time.Second {
			delete(C.data, key)
		}
	}
}