package lru

import (
	"container/list"
	"sync"
)

type Cacher[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Delete(key K)
}

type LRUCache[K comparable, V any] struct {
	capacity int
	// cache is a map that holds the key and a pointer to the list element
	// This allows O(1) access to the elements in the cache
	cache map[K]*list.Element
	// list is a doubly linked list that holds the values in the order of their usage
	list *list.List
	mu   sync.Mutex
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
	}
}

// Get retrieves the value for a key from the cache.
// it checks if the key exists in the cache
// If it does, move the element to the front of the list as recently used and return its value
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves it to the front of the
// list as recently used. If the cache is at capacity, it removes the least recently used
// item before adding the new item.
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If the key already exists, update its value and move it to the front of the list
	// If it doesn't exist, check if the cache is at capacity
	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		elem.Value = entry[K, V]{key: key, value: value}
		return
	}

	if c.list.Len() >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			delete(c.cache, oldest.Value.(entry[K, V]).key)
		}
	}
	// Add the new key-value pair to the cache and the front of the list
	elem := c.list.PushFront(entry[K, V]{key: key, value: value})
	c.cache[key] = elem
}

// Delete deletes a key-value pair from the cache.
func (c *LRUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cache[key]; found {
		c.list.Remove(elem)
		delete(c.cache, key)
	}
}
