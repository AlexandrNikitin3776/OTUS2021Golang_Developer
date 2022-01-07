package hw04lrucache

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

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
