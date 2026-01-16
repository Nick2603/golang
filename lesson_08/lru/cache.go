package lru

import (
	"container/list"
	"errors"
)

var ErrInvalidCapacity = errors.New("capacity must be positive")

type Cache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type cacheItem struct {
	key   string
	value string
}

type cache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

func NewCache(capacity int) (Cache, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &cache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}, nil
}

func (c *cache) Get(key string) (string, bool) {
	if element, exists := c.cache[key]; exists {
		c.list.MoveToFront(element)
		item := element.Value.(*cacheItem)
		return item.value, true
	}

	return "", false
}

func (c *cache) Put(key, value string) {
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
