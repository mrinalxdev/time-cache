package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mu        sync.RWMutex
	ttl       time.Duration
	cache     map[string]*list.Element
	lruList   *list.List
	onEvicted func(string, interface{})
}

// Utilities functions so that I can use it elsewhere in the program
type entry struct {
	key   string
	value interface{}
	ttl   time.Time
}

// Removes a key-value pair from the cache.
func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		delete(c.cache, key)
		c.lruList.Remove(elem)
		c.onEvicted(key, elem.Value.(*entry).value)
	}
}

func (c *Cache) remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Remove(key)
}

// NewCache gives new cache instance
// initializes the cache with time-to-live and an eviction callback
// ttl : time duration for which a cache item is valid
// cache : a map where the key is a string and the value is a pointer to the element in the LRU list.
// this allows quick access to cache items.
func NewCache(ttl time.Duration, onEvicted func(string, interface{})) *Cache {
	return &Cache{
		ttl:       ttl,
		cache:     make(map[string]*list.Element),
		lruList:   list.New(),
		onEvicted: onEvicted,
	}
}

// :TODO retrieve a new cache instance
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	// Move accessed element to front of LRU list.
	c.lruList.MoveToFront(elem)
	return elem.Value.(*entry).value, true
}

// :TODO retieve a value and its ttl from the cache.
func (c *Cache) GetWithTTL(key string) (interface{}, time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.cache[key]
	if !ok {
		return nil, 0, false
	}
	ttl := elem.Value.(*entry).ttl.Sub(time.Now())
	return elem.Value.(*entry).value, ttl, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		elem.Value.(*entry).value = value
		c.lruList.MoveToFront(elem)
		return
	}

	entry := &entry{
		key:   key,
		value: value,
		ttl:   time.Now().Add(c.ttl),
	}

	c.cache[key] = c.lruList.PushFront(entry)

	time.AfterFunc(c.ttl, func() {
		c.remove(key)
	})
}

func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &entry{
		key:   key,
		value: value,
		ttl:   time.Now().Add(ttl),
	}
	c.cache[key] = c.lruList.PushFront(entry)
}

// removes the oldest item from the cache
func (c *Cache) DeleteOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lruList.Len() == 0 {
		return
	}

	oldest := c.lruList.Back()
	c.lruList.Remove(oldest)
	delete(c.cache, oldest.Value.(*entry).key)
	c.onEvicted(oldest.Value.(*entry).key, oldest.Value.(*entry).value)
}

// it returns the nimber of items in the cache
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lruList.Len()
}

// removes all the items from the cache
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, elem := range c.cache {
		c.onEvicted(elem.Value.(*entry).key, elem.Value.(*entry).value)
	}

	c.cache = make(map[string]*list.Element)
	c.lruList.Init()
}

// returns a slice of all key in the cache
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, c.lruList.Len())
	for elem := c.lruList.Front(); elem != nil; elem = elem.Next() {
		keys = append(keys, elem.Value.(*entry).key)
	}

	return keys
}

// returns a slice of all values in the cache.
func (c *Cache) Values() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	values := make([]interface{}, 0, c.lruList.Len())
	for elem := c.lruList.Front(); elem != nil; elem = elem.Next() {
		values = append(values, elem.Value.(*entry).value)
	}

	return values
}

// return the time-to-love for the cache
func (c *Cache) TTL() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.ttl
}

// updates the ttl for a key
func (c *Cache) UpdateTTL(key string, ttl time.Duration) {
	c.mu.RLock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		elem.Value.(*entry).ttl = time.Now().Add(ttl)
	}
}

// checks if a key exists in the cache
func (c *Cache) CheckExists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.cache[key]
	return ok
}

// returns the expiration time for a key
func (c *Cache) GetExpirationTime(key string) (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.cache[key]
	if !ok {
		return time.Time{}, false
	}

	return elem.Value.(*entry).ttl, true
}

func (c *Cache) DeleteBefore(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, elem := range c.cache {
		if elem.Value.(*entry).ttl.Before(t) {
			c.remove(elem.Value.(*entry).key)
		}
	}
}

func (c *Cache) DeleteAfter(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, elem := range c.cache {
		if elem.Value.(*entry).ttl.After(t) {
			c.remove(elem.Value.(*entry).key)
		}
	}
}

// returns cache statistics
func (c *Cache) Stats() (int, int, int) {

}
