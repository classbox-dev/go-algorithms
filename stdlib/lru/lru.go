package lru

import (
	"container/list"
)

type pair[K comparable, V comparable] struct {
	key   K
	value V
}

// Cache is a fixed-size LRU cache for integer keys and values
type Cache[K comparable, V comparable] struct {
	cap int
	cl  *list.List
	cm  map[K]*list.Element
}

// New returns an initialised LRU cache of the given capacity
func New[K comparable, V comparable](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		cap: capacity,
		cl:  list.New(),
		cm:  make(map[K]*list.Element),
	}
}

// Get returns a cached value for the given key and a flag indicating whether the key is in the cache.
func (cache *Cache[K, V]) Get(key K) (V, bool) {
	node, ok := cache.cm[key]
	if !ok {
		var zero V
		return zero, false
	}
	cache.cl.MoveToFront(node)
	item := node.Value.(pair[K, V])
	return item.value, true
}

// Put inserts or updates the value for the given key.
// If the cache capacity is reached, Put removes the least recently used item before inserting a new one.
// A key is considered "used" every time it is fetched with Get or its value is updated with Put.
func (cache *Cache[K, V]) Put(key K, value V) {
	node, ok := cache.cm[key]
	if ok {
		item := node.Value.(pair[K, V])
		if item.value != value {
			item.value = value
			node.Value = item
		}
		cache.cl.MoveToFront(node)
		return
	}

	if len(cache.cm) >= cache.cap {
		node = cache.cl.Back()
		item := node.Value.(pair[K, V])
		delete(cache.cm, item.key)
		cache.cl.Remove(node)
	}
	item := pair[K, V]{key, value}
	cache.cm[key] = cache.cl.PushFront(item)
}
