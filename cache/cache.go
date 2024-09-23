package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mu sync.RWMutex
	ttl time.Duration
	cache map[string]*list.Element
	lruList *list.List
	onEvicted func(string, interface{})
}

type entry struct {
	key string
	value interface {}
	ttl time.Time
}

// :TODO retrieve a new cache instance
func NewCache(ttl time.Duration, onEvicted func(string, interface{})) *Cache {
	return &Cache {
		ttl : ttl,
		cache : make(map[string]*list.Element),
		lruList : list.New(),
		onEvicted: onEvicted,
	}
}
func (c *Cache) Get(key string) (interface {}, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.lruList.MoveToFront(elem)
	return elem.Value.(*entry).value, true
}


