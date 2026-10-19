package repositories

import (
	"container/list"
	"sync"
)

// entry stores the actual key-value data pair
type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity  int
	mu        sync.Mutex
	evictList *list.List
	items     map[K]*list.Element
}

// initialize new LRUCache with a fixed capacity
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity:  capacity,
		evictList: list.New(),
		items:     make(map[K]*list.Element),
	}
}

func (c *LRUCache[K, V]) Get(key K) (value V, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exist := c.items[key]; exist {
		c.evictList.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	return value, false
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.evictList.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	kv := &entry[K, V]{key: key, value: value}
	elem := c.evictList.PushFront(kv)
	c.items[key] = elem

	if c.evictList.Len() > c.capacity {
		c.evictOldest()
	}
}

func (c *LRUCache[K, V]) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.evictList.Remove(elem)
		kv := elem.Value.(*entry[K, V])
		delete(c.items, kv.key)
	}
}
