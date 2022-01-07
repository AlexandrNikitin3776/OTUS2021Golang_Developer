package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lock     sync.Mutex
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

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.lock.Lock()
	ok := cache.SetItem(key, value)
	cache.lock.Unlock()
	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.lock.Lock()
	value, ok := cache.GetItem(key)
	cache.lock.Unlock()
	return value, ok
}

func (cache *lruCache) Clear() {
	cache.lock.Lock()
	cache.ClearCache()
	cache.lock.Unlock()
}

func (cache *lruCache) SetItem(key Key, value interface{}) bool {
	if item, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(item)
		item.Value.(*cacheItem).value = value
		return true
	}

	newCacheItem := &cacheItem{key, value}
	newQueueItem := cache.queue.PushFront(newCacheItem)
	cache.items[key] = newQueueItem

	if cache.queue.Len() > cache.capacity {
		lastQueueItem := cache.queue.Back()
		cache.queue.Remove(lastQueueItem)
		delete(cache.items, lastQueueItem.Value.(*cacheItem).key)
	}
	return false
}

func (cache *lruCache) GetItem(key Key) (interface{}, bool) {
	if item, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) ClearCache() {
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}
