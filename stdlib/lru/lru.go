package lru

import (
	"container/list"
)

type pair struct {
	key, value int
}

// Cache is a fixed-size LRU cache for integer keys and values
type Cache struct {
	cap int
	cl  *list.List
	cm  map[int]*list.Element
}

// New returns an initialised LRU cache of the given capacity
func New(capacity int) *Cache {
	return &Cache{
		cap: capacity,
		cl:  list.New(),
		cm:  make(map[int]*list.Element),
	}
}

// Get returns a cached value for the given key and a flag indicating whether the key is in the cache.
func (cache *Cache) Get(key int) (int, bool) {
	node, ok := cache.cm[key]
	if !ok {
		return 0, false
	}
	cache.cl.MoveToFront(node)
	item := node.Value.(pair)
	return item.value, true
}

// Put inserts or updates the value for the given key.
// If the cache capacity is reached, Put removes the least recently used item before inserting a new one.
// A key is considered "used" every time it is fetched with Get or its value is updated with Put.
func (cache *Cache) Put(key int, value int) {
	node, ok := cache.cm[key]
	if ok {
		item := node.Value.(pair)
		if item.value != value {
			item.value = value
			node.Value = item
		}
		cache.cl.MoveToFront(node)
		return
	}

	if len(cache.cm) >= cache.cap {
		node = cache.cl.Back()
		item := node.Value.(pair)
		delete(cache.cm, item.key)
		cache.cl.Remove(node)
	}
	item := pair{key, value}
	cache.cm[key] = cache.cl.PushFront(item)
}
