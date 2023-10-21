package memcache

import (
	"sync"
	"time"
)

// Cache is a generic in-memory key-value thread-safe* cache.
//
// *Due to the generic nature of the cache it is possible to store types that
// are mutatable by reference which is not thread-safe. Instead of applying a
// stricter type constraint on K or V to prevent this, it is left up to the user
// to decide the nature of their cache.
type Cache[K comparable, V any] struct {
	mu                 sync.RWMutex
	store              map[K]Item[K, V]
	passiveExpiration  bool
	expirationInterval time.Duration
	expirer            ExpirerFunc[K, V]
	closeCh            chan struct{}
	closed             bool
}

// Open a new in-memory key-value cache.
func Open[K comparable, V any](options ...Option[K, V]) (*Cache[K, V], error) {
	cache := &Cache[K, V]{
		mu:      sync.RWMutex{},
		store:   map[K]Item[K, V]{},
		closeCh: make(chan struct{}),
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
		go cache.runActiveExpirer()
	}

	return cache, nil
}

// Set permanent key to hold value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = Item[K, V]{Value: value}
}

// SetEx key to hold value in the cache and set key to timeout after the
// provided ttl.
func (c *Cache[K, V]) SetEx(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expireAt := time.Now().Add(ttl)
	c.store[key] = Item[K, V]{Value: value, ExpireAt: &expireAt}
}

// Get returns the value associated with the provided key if it exists, or false
// if it does not.
//
// If the cache was opened with [WithPassiveExpiration] and the requested key
// is expired, it will be deleted from the cache and false will be returned.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	item, ok := c.store[key]
	c.mu.RUnlock()

	if ok && c.passiveExpiration && item.IsExpired() {
		c.Delete(key)
		var v V
		return v, false
	}

	return item.Value, ok
}

// Has returns true if the provided key exists in the cache.
//
// If the cache was opened with [WithPassiveExpiration] and the requested key
// is expired, it will be deleted from the cache and false will be returned.
func (c *Cache[K, V]) Has(key K) bool {
	_, ok := c.Get(key)
	return ok
}

// Delete provided keys from the cache.
func (c *Cache[K, V]) Delete(keys ...K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.store, key)
	}
}

// Flush the cache, deleting all keys.
func (c *Cache[K, V]) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.store)
}

// Size returns the number of items currently in the cache.
func (c *Cache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.store)
}

// Keys returns a slice of all keys currently in the cache.
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.store))
	for key := range c.store {
		keys = append(keys, key)
	}

	return keys
}

// Close the cache, stopping all running goroutines. Should be called when the
// cache is no longer needed.
func (c *Cache[K, V]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
	close(c.closeCh)
}

func (c *Cache[K, V]) runActiveExpirer() {
	ticker := time.NewTicker(c.expirationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeCh:
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
