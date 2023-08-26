package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	m        sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.m.Lock()
	defer lc.m.Unlock()
	return lc.set(key, value)
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.m.Lock()
	defer lc.m.Unlock()
	return lc.get(key)
}

func (lc *lruCache) set(key Key, value interface{}) bool {
	item, ok := lc.items[key]

	switch ok {
	case true:
		item.Value.(*cacheItem).value = value
		lc.queue.MoveToFront(item)
	case false:
		if lc.capacity == lc.queue.Len() {
			last := lc.queue.Back()
			item := last.Value.(*cacheItem)
			lc.queue.Remove(last)
			delete(lc.items, item.key)
		}
		ci := &cacheItem{
			key:   key,
			value: value,
		}
		lc.items[ci.key] = lc.queue.PushFront(ci)
	}
	return ok
}

func (lc *lruCache) get(key Key) (interface{}, bool) {
	item, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(item)
		item := lc.queue.Front().Value
		return item.(*cacheItem).value, ok
	}
	return nil, ok
}

func (lc *lruCache) Clear() {
	lc.m.Lock()
	defer lc.m.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
