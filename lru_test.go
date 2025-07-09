package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache_PutAndGet(t *testing.T) {
	cache := NewLRUCache[int, string](2)

	cache.Put(1, "one")
	cache.Put(2, "two")

	val, ok := cache.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "one", val)

	val, ok = cache.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "two", val)

	val, ok = cache.Get(3)
	assert.False(t, ok)
}

func TestLRUCache_EvictsLeastRecentlyUsed(t *testing.T) {
	cache := NewLRUCache[int, string](2)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Get(1)          // 1 is most recently used
	cache.Put(3, "three") // should evict 2

	if _, ok := cache.Get(2); ok {
		t.Errorf("expected key 2 to be evicted")
	}
	if val, ok := cache.Get(1); !ok || val != "one" {
		t.Errorf("expected to get 'one', got '%v', ok=%v", val, ok)
	}
	if val, ok := cache.Get(3); !ok || val != "three" {
		t.Errorf("expected to get 'three', got '%v', ok=%v", val, ok)
	}
}

func TestLRUCache_UpdateValue(t *testing.T) {
	cache := NewLRUCache[int, string](2)

	cache.Put(1, "one")
	cache.Put(1, "uno")

	if val, ok := cache.Get(1); !ok || val != "uno" {
		t.Errorf("expected to get 'uno', got '%v', ok=%v", val, ok)
	}
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache[int, string](2)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Delete(1)

	if _, ok := cache.Get(1); ok {
		t.Errorf("expected key 1 to be deleted")
	}
	if val, ok := cache.Get(2); !ok || val != "two" {
		t.Errorf("expected to get 'two', got '%v', ok=%v", val, ok)
	}
}

func TestLRUCache_GetNonExistent(t *testing.T) {
	cache := NewLRUCache[int, string](2)

	if val, ok := cache.Get(42); ok {
		t.Errorf("expected to get false for non-existent key, got '%v', ok=%v", val, ok)
	}
}
