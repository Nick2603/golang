package lru

import (
	"testing"
)

func TestLruCache_PutAndGet(t *testing.T) {
	cache := NewLruCache(2)

	cache.Put("key1", "value1")
	if value, ok := cache.Get("key1"); !ok || value != "value1" {
		t.Errorf("Expected 'value1', got '%s'", value)
	}

	cache.Put("key2", "value2")
	if value, ok := cache.Get("key2"); !ok || value != "value2" {
		t.Errorf("Expected 'value2', got '%s'", value)
	}

	if value, ok := cache.Get("key1"); !ok || value != "value1" {
		t.Errorf("Expected 'value1' after adding 'key2', got '%s'", value)
	}
}

func TestLruCache_ReplaceLeastRecentlyUsed(t *testing.T) {
	cache := NewLruCache(2)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	cache.Put("key3", "value3")

	if _, ok := cache.Get("key1"); ok {
		t.Errorf("Expected 'key1' to be evicted")
	}

	if value, ok := cache.Get("key2"); !ok || value != "value2" {
		t.Errorf("Expected 'value2', got '%s'", value)
	}

	if value, ok := cache.Get("key3"); !ok || value != "value3" {
		t.Errorf("Expected 'value3', got '%s'", value)
	}
}

func TestLruCache_UpdateKey(t *testing.T) {
	cache := NewLruCache(2)

	cache.Put("key1", "value1")
	cache.Put("key1", "updated_value1")

	if value, ok := cache.Get("key1"); !ok || value != "updated_value1" {
		t.Errorf("Expected 'updated_value1', got '%s'", value)
	}
}

func TestLruCache_LRUBehavior(t *testing.T) {
	cache := NewLruCache(2)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Get("key1")
	cache.Put("key3", "value3")

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("Expected 'key2' to be evicted")
	}

	if value, ok := cache.Get("key1"); !ok || value != "value1" {
		t.Errorf("Expected 'value1', got '%s'", value)
	}

	if value, ok := cache.Get("key3"); !ok || value != "value3" {
		t.Errorf("Expected 'value3', got '%s'", value)
	}
}
