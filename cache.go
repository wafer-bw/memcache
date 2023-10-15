package memcache

import (
	"context"
	"sync"

	"github.com/wafer-bw/memcache/internal/record"
)

// TODO: decide how to handle options that need generics such as:
// WithExpirer[K, V]
// WithEvicter[K, V]
// WithOnEvict[K, V]
//
// The base cache options can reasonably be generic but the record options
// are likely best left as concrete types. This would mean controlling records
// generically must be done on separate cache methods.

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

	recordConfig := RecordConfig{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&recordConfig)
	}

	c.store[key] = record.Record[V]{
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

func (c *Cache[K, V]) Length() int {
	// TODO: rename to Size().

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
