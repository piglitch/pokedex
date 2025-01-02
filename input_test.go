package main

import (
	"sync"
	"testing"
	"time"

	"github.com/piglitch/pokedexcli/pokecache"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := pokecache.NewCache(10 * time.Second)

	// Add an item to the cache
	key := "testKey"
	val := []byte("testValue")
	cache.Add(key, val)

	// Retrieve the item
	retrievedVal, exists := cache.Get(key)
	if !exists {
		t.Fatalf("expected key %s to exist", key)
	}
	if string(retrievedVal) != string(val) {
		t.Errorf("expected value %s, got %s", string(val), string(retrievedVal))
	}
}

func TestCache_ItemExpiration(t *testing.T) {
	cache := pokecache.NewCache(1 * time.Second)

	// Add an item to the cache
	key := "testKey"
	val := []byte("testValue")
	cache.Add(key, val)

	// Wait for the item to expire
	time.Sleep(6 * time.Second)

	// Try to retrieve the item
	_, exists := cache.Get(key)
	if exists {
		t.Fatalf("expected key %s to be expired and removed", key)
	}
}

func TestCache_ReapLoop(t *testing.T) {
	cache := pokecache.NewCache(1 * time.Second)

	// Add multiple items with different creation times
	cache.Add("key1", []byte("value1"))
	time.Sleep(3 * time.Second)
	cache.Add("key2", []byte("value2"))

	// Wait for reaping to remove old items
	time.Sleep(3 * time.Second)

	// Check which items remain
	_, exists1 := cache.Get("key1")
	if exists1 {
		t.Errorf("expected key1 to be reaped, but it still exists")
	}
	_, exists2 := cache.Get("key2")
	if !exists2 {
		t.Errorf("expected key2 to still exist, but it was reaped")
	}
}

func TestCache_Concurrency(t *testing.T) {
	cache := pokecache.NewCache(10 * time.Second)

	// Use multiple goroutines to test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('A' + i))
			val := []byte{byte(i)}
			cache.Add(key, val)
			if _, exists := cache.Get(key); !exists {
				t.Errorf("key %s was not found after being added", key)
			}
		}(i)
	}
	wg.Wait()

	// Ensure all keys are present
	for i := 0; i < 100; i++ {
		key := string(rune('A' + i))
		if _, exists := cache.Get(key); !exists {
			t.Errorf("key %s was not found", key)
		}
	}
}

func TestCache_EmptyGet(t *testing.T) {
	cache := pokecache.NewCache(10 * time.Second)

	// Try to get a key that doesn't exist
	_, exists := cache.Get("nonExistentKey")
	if exists {
		t.Errorf("expected key to not exist, but it was found")
	}
}
