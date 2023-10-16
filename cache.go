package memcache

import (
	"context"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	mu                sync.RWMutex
	store             map[K]Item[K, V]
	passiveExpiration bool

	expirationInterval time.Duration
	expirer            ExpirerFunc[K, V]
}

func New[K comparable, V any](ctx context.Context, options ...CacheOption[K, V]) (*Cache[K, V], error) {
	cache := &Cache[K, V]{
		mu:    sync.RWMutex{},
		store: map[K]Item[K, V]{},
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(cache); err != nil {
			return nil, err
		}
	}

	if cache.expirer != nil && cache.expirationInterval > 0 {
		go cache.runExpirer(ctx)
	}

	return cache, nil
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = Item[K, V]{Value: value}
}

func (c *Cache[K, V]) SetEx(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expireAt := time.Now().Add(ttl)
	c.store[key] = Item[K, V]{Value: value, ExpireAt: &expireAt}
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

func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}

	return keys
}

func (c *Cache[K, V]) runExpirer(ctx context.Context) {
	ticker := time.NewTicker(c.expirationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// TODO: this allows the context to be closed/cancelled but does not
			//       stop the cache from being used. Should consider a solution
			//       to this problem.
			return
		case <-ticker.C:
			c.mu.Lock()
			c.expirer(c.store)
			c.mu.Unlock()
		}
	}
}

// TODO: Add Items() that returns a shallow copy of c.store?
//       Naming may need to distinguish from a potential method that
//       gives access to the real c.store.
