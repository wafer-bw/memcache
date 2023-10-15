package memcache

import (
	"context"
	"sync"

	"github.com/wafer-bw/memcache/internal/record"
)

type Cache[K comparable, V any] struct {
	mu                sync.RWMutex
	store             map[K]record.Record[V]
	passiveExpiration bool

	// TODO: add active expiration support.
	// TODO: add active eviction support.
}

func New[K comparable, V any](ctx context.Context, options ...CacheConfigOption) (*Cache[K, V], error) {
	config := CacheConfig{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&config); err != nil {
			return nil, err
		}
	}

	cache := &Cache[K, V]{
		mu:                sync.RWMutex{},
		store:             map[K]record.Record[V]{},
		passiveExpiration: config.passiveExpiration,
	}

	return cache, nil
}

func (c *Cache[K, V]) Set(key K, value V, options ...RecordConfigOption) {
	c.mu.Lock()
	defer c.mu.Unlock()

	valueConfig := RecordConfig{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&valueConfig)
	}

	c.store[key] = record.Record[V]{
		Value:    value,
		ExpireAt: valueConfig.expireAt,
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

func (c *Cache[K, V]) Length() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.store)
}

func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}

	return keys
}
