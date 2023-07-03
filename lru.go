package golru

import (
	"fmt"
	"sync"

	"github.com/thepatrik/golru/internal/lru"
)

type cacheable interface {
	any
}

// OnEvicted is a callback function that is called when an item is evicted from the cache.
type OnEvicted[K comparable, V cacheable] func(key K, value V)

// LRUCache is a thread safe LRU cache.
type LRUCache[K comparable, V any] struct {
	lru       *lru.LRUCache[K, V]
	lock      sync.RWMutex
	onEvicted OnEvicted[K, V]
}

// Option type.
type Option func(*Config)

// Config for LRUCache.
type Config struct {
	MaxEntries int
}

// WithMaxEntries sets max entries before eviction. Default is 0 (no limit).
func WithMaxEntries(maxEntries int) Option {
	return func(cfg *Config) {
		cfg.MaxEntries = maxEntries
	}
}

// New returns a new thread safe LRU cache.
func New[K comparable, V any](options ...Option) (*LRUCache[K, V], error) {
	cfg := &Config{
		MaxEntries: 0,
	}
	for _, option := range options {
		option(cfg)
	}

	if cfg.MaxEntries < 0 {
		return nil, fmt.Errorf("max entries cannot be negative")
	}

	lru := lru.New[K, V](cfg.MaxEntries)

	cache := &LRUCache[K, V]{
		lru: lru,
	}

	return cache, nil
}

// SetOnEvicted sets callback func for evictions.
func (cache *LRUCache[K, V]) SetOnEvicted(onEvicted OnEvicted[K, V]) {
	cache.onEvicted = onEvicted
}

// Contains checks if a key is in the cache.
func (cache *LRUCache[K, V]) Contains(key K) bool {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	return cache.lru.Contains(key)
}

// Get returns the value of a key in the cache and a bool indicating if the cache item exists.
func (cache *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	return cache.lru.Get(key)
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (cache *LRUCache[K, V]) Keys() []K {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	return cache.lru.Keys()
}

// Len returns the length of the cache.
func (cache *LRUCache[K, V]) Len() int {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	return cache.lru.Len()
}

// Peek returns the value of a key in the cache without updating resentness of the cached item.
func (cache *LRUCache[K, V]) Peek(key K) (value V, ok bool) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	return cache.lru.Peek(key)
}

// Put puts an item to the cache.
func (cache *LRUCache[K, V]) Put(key K, val V) {
	ok, evicted := func() (bool, *lru.Item[K, V]) {
		cache.lock.Lock()
		defer cache.lock.Unlock()
		return cache.lru.Put(key, val)
	}()

	if cache.onEvicted != nil && ok {
		cache.onEvicted(evicted.Key, evicted.Val)
	}
}

// Remove removes an item from the cache and returns a bool indicating if the removal was successful.
func (cache *LRUCache[K, V]) Remove(key K) {
	ok, evicted := func() (bool, *lru.Item[K, V]) {
		cache.lock.Lock()
		defer cache.lock.Unlock()
		return cache.lru.Remove(key)
	}()

	if cache.onEvicted != nil && ok {
		cache.onEvicted(evicted.Key, evicted.Val)
	}
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (cache *LRUCache[K, V]) Values() []V {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	return cache.lru.Values()
}
