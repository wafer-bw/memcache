package memcache

import (
	"context"
	"sync"
	"time"

	"github.com/wafer-bw/memcache/internal/record"
)

type Cache[K comparable, V any] struct {
	mu                 sync.RWMutex
	store              map[K]record.Record[V]
	expireOnGet        bool
	expirationInterval time.Duration

	// TODO: add eviction support.
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
		mu:                 sync.RWMutex{},
		store:              map[K]record.Record[V]{},
		expireOnGet:        config.expireOnGet,
		expirationInterval: config.expirationInterval,
	}

	if cache.expirationInterval > 0 {
		go cache.expirer(ctx)
	}

	return cache, nil
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

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	r, ok := c.store[key]
	if !ok {
		return r.Value, ok
	}
	c.mu.RUnlock()

	if c.expireOnGet && r.IsExpired() {
		c.Delete(key)
		var v V
		return v, false
	}

	return r.Value, ok
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
	expireOnGet        bool
	expirationInterval time.Duration
}

type CacheConfigOption func(*CacheConfig) error

func WithExpireOnGet() CacheConfigOption {
	return func(config *CacheConfig) error {
		config.expireOnGet = true
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

func (c *Cache[K, V]) expirer(ctx context.Context) {
	// TODO: add unit tests for this.
	ticker := time.NewTicker(c.expirationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			// TODO: employ a more efficient algorithm for expiring records.
			// https://redis.io/commands/expire/#:~:text=How%20Redis%20expires%20keys
			for k, r := range c.store {
				if r.IsExpired() {
					delete(c.store, k)
				}
			}
			c.mu.Unlock()
		}
	}
}
