package memcache

import (
	"sync"
	"time"

	"github.com/wafer-bw/memcache/internal/record"
)

type Cache[K comparable, V any] struct {
	mu                 sync.RWMutex
	store              map[K]record.Record[V]
	expirationInterval time.Duration
	evictionInterval   time.Duration
}

func New[K comparable, V any](options ...CacheConfigOption) (*Cache[K, V], error) {
	config := CacheConfig{}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&config); err != nil {
			return nil, err
		}
	}

	return &Cache[K, V]{
		mu:                 sync.RWMutex{},
		store:              map[K]record.Record[V]{},
		evictionInterval:   config.evictionInterval,
		expirationInterval: config.expirationInterval,
	}, nil
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	r, ok := c.store[key]

	return r.Value, ok
}

func (c *Cache[K, V]) Set(key K, value V, options ...ValueConfigOption) {
	c.mu.Lock()
	defer c.mu.Unlock()

	valueConfig := ValueConfig{}
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

type CacheConfig struct {
	// mode             uint8
	evictionInterval   time.Duration
	expirationInterval time.Duration
}

type CacheConfigOption func(*CacheConfig) error

func WithEvictionInterval(i time.Duration) CacheConfigOption {
	return func(config *CacheConfig) error {
		config.evictionInterval = i
		return nil
	}
}

func WithExpirationInterval(i time.Duration) CacheConfigOption {
	return func(config *CacheConfig) error {
		config.expirationInterval = i
		return nil
	}
}

type ValueConfig struct {
	expireAt *time.Time
}

type ValueConfigOption func(*ValueConfig)

func WithTTL(d time.Duration) ValueConfigOption {
	return func(config *ValueConfig) {
		expireAt := time.Now().Add(d)
		config.expireAt = &expireAt
	}
}
