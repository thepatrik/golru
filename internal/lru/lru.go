package lru

type cacheable interface {
	any
}

type Item[K comparable, V cacheable] struct {
	Key        K
	Val        V
	next, prev *Item[K, V]
}

// LRUCache is an internal thread unsafe LRU cache.
type LRUCache[K comparable, V cacheable] struct {
	first, last *Item[K, V]
	items       map[K]*Item[K, V]
	// MaxEntries is the maximum number of cache entries before eviction.
	MaxEntries int
}

// New returns a new LRU cache.
func New[K comparable, V cacheable](maxEntries int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		MaxEntries: maxEntries,
		items:      make(map[K]*Item[K, V], 0),
	}
}

// Contains checks if a key is in the cache.
func (cache *LRUCache[K, V]) Contains(key K) bool {
	_, ok := cache.items[key]
	return ok
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (cache *LRUCache[K, V]) Keys() []K {
	keys := make([]K, len(cache.items))
	item := cache.last

	for i := 0; i < len(cache.items); i++ {
		keys[i] = item.Key
		item = item.prev
	}

	return keys
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (cache *LRUCache[K, V]) Values() []V {
	values := make([]V, len(cache.items))
	item := cache.last

	for i := 0; i < len(cache.items); i++ {
		values[i] = item.Val
		item = item.prev
	}

	return values
}

// Get returns the value of a key in the cache and a bool indicating if the cache item exists.
func (cache *LRUCache[K, V]) Get(key K) (V, bool) {
	item, ok := cache.items[key]
	if !ok {
		var val V
		return val, ok
	}

	if cache.first.Key == key {
		return item.Val, true
	}

	cache.moveToFront(item)

	return item.Val, true
}

// Len returns the length of the cache.
func (cache *LRUCache[K, V]) Len() int {
	return len(cache.items)
}

// Peek returns the value of a key in the cache without updating resentness of the cached item.
func (cache *LRUCache[K, V]) Peek(key K) (val V, ok bool) {
	item, ok := cache.items[key]
	if ok {
		val = item.Val
	}

	return val, ok
}

// Put puts an item to the cache.
func (cache *LRUCache[K, V]) Put(key K, value V) (bool, *Item[K, V]) {
	cachedItem, ok := cache.items[key]
	if ok {
		cachedItem.Val = value
		cache.items[key] = cachedItem
		return false, nil
	}

	var evictedItem *Item[K, V]

	if cache.MaxEntries > 0 && len(cache.items) >= cache.MaxEntries {
		last := cache.last
		if last != nil {
			cache.remove(last)
			evictedItem = last
		}
	}

	item := &Item[K, V]{
		Key:  key,
		Val:  value,
		next: cache.first,
	}

	cache.moveToFront(item)
	if evictedItem != nil {
		return true, evictedItem
	}

	return false, nil
}

// Remove removes an item from the cache and returns a bool indicating if the removal was successful.
func (cache *LRUCache[K, V]) Remove(key K) (bool, *Item[K, V]) {
	item, ok := cache.items[key]
	if !ok {
		return false, nil
	}

	cache.remove(item)
	return true, item
}

func (cache *LRUCache[K, V]) moveToFront(item *Item[K, V]) {
	prev := item.prev
	next := item.next
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}

	first := cache.first
	if first != nil {
		first.prev = item
	}

	last := cache.last

	if last == nil {
		cache.last = item
	} else if last.Key == item.Key {
		prev = last.prev
		if prev != nil {
			prev.next = nil
		}
		cache.last = prev
	}

	item.prev = nil
	item.next = first
	cache.first = item

	cache.items[item.Key] = item
}

func (cache *LRUCache[K, V]) remove(item *Item[K, V]) {
	prev := item.prev
	next := item.next
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}

	first := cache.first
	if first != nil && first.Key == item.Key {
		cache.first = item.next
	}

	last := cache.last
	if last != nil && last.Key == item.Key {
		cache.last = item.prev
	}

	delete(cache.items, item.Key)
}
