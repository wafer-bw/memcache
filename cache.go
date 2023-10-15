package memcache

import (
	"context"
	"sync"
)

// TODO: decide how to handle item options that need generics such as:
// WithOnEvict[K, V]

type Cache[K comparable, V any] struct {
	mu                sync.RWMutex
	store             map[K]Item[V]
	passiveExpiration bool
}

func New[K comparable, V any](ctx context.Context, options ...CacheOption[K, V]) (*Cache[K, V], error) {
	cache := &Cache[K, V]{
		mu:    sync.RWMutex{},
		store: map[K]Item[V]{},
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(cache); err != nil {
			return nil, err
		}
	}

	return cache, nil
}

func (c *Cache[K, V]) Set(key K, value V, options ...ItemOption) {
	c.mu.Lock()
	defer c.mu.Unlock()

	recordConfig := ItemConfig{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&recordConfig)
	}

	c.store[key] = Item[V]{
		Value:    value,
		ExpireAt: recordConfig.expireAt,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	r, ok := c.store[key]
	c.mu.RUnlock()

	if ok && c.passiveExpiration && r.IsExpired() {
		c.Delete(key)
		var v V
		return v, false
	}

	return r.Value, ok
}

func (c *Cache[K, V]) Has(key K) bool {
	_, ok := c.Get(key)
	return ok
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}

func (c *Cache[K, V]) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.store)
}

func (c *Cache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.store)
}

// TODO: add SizeBytes() that returns the size of the cache in bytes.

func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}

	return keys
}

// TODO: Add Items() that returns a shallow copy of c.store.
