package hw04lrucache

import "sync"

type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type lruCache struct {
	items    map[string]*cacheMapElement
	capacity int
	queue    List
	sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		items:    map[string]*cacheMapElement{},
		capacity: capacity,
		queue:    NewList(),
	}
}

type cacheMapElement struct {
	el    *ListItem
	Value interface{}
}

func (c *lruCache) Set(key string, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	v, ok := c.items[key]
	if !ok {
		el := c.queue.PushFront(key)
		c.items[key] = &cacheMapElement{
			el:    el,
			Value: value,
		}

		if c.queue.Len() > c.capacity {
			backEl := c.queue.Back()
			backElementKey := backEl.Value
			c.queue.Remove(backEl)
			delete(c.items, backElementKey.(string))
		}
	} else {
		v.Value = value
		c.queue.MoveToFront(v.el)
		return true
	}
	return false
}

func (c *lruCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(v.el)

	return v.Value, true
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = map[string]*cacheMapElement{}
}
