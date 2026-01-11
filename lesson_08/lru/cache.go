package lru

import (
	"container/list"
)

type LruCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type cacheItem struct {
	key   string
	value string
}

type lruCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

func NewLruCache(capacity int) LruCache {
	if capacity <= 0 {
		panic("capacity must be positive")
	}

	return &lruCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *lruCache) Get(key string) (string, bool) {
	if element, exists := c.cache[key]; exists {
		c.list.MoveToFront(element)
		item := element.Value.(*cacheItem)
		return item.value, true
	}

	return "", false
}

func (c *lruCache) Put(key, value string) {
	if element, exists := c.cache[key]; exists {
		c.list.MoveToFront(element)
		item := element.Value.(*cacheItem)
		item.value = value

		return
	}

	item := &cacheItem{
		key:   key,
		value: value,
	}

	element := c.list.PushFront(item)

	c.cache[key] = element

	if c.list.Len() > c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			oldestItem := oldest.Value.(*cacheItem)
			delete(c.cache, oldestItem.key)
		}
	}
}
