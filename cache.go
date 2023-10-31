package memcache

import (
	"time"

	"github.com/wafer-bw/memcache/internal/closer"
)

// Cache is a generic in-memory key-value thread-safe* cache.
//
// *Due to the generic nature of the cache it is possible to store types that
// are mutatable by reference which is not thread-safe. Instead of applying a
// stricter type constraint on K & V to prevent this, it is left up to the user
// to decide the nature of their cache.
type Cache[K comparable, V any] struct {
	store  storer[K, V]
	closer *closer.Closer

	passiveExpiration        bool
	activeExpirationInterval time.Duration
}

// Open a new in-memory key-value cache.
func Open[K comparable, V any](options ...Option[K, V]) (*Cache[K, V], error) {
	cache := &Cache[K, V]{
		store:  newNoEvictStore[K, V](),
		closer: closer.New(),
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(cache); err != nil {
			return nil, err
		}
	}

	if cache.activeExpirationInterval > 0 {
		go cache.runActiveExpirer()
	}

	return cache, nil
}

// Set permanent key to hold value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.store.Set(key, Item[K, V]{Value: value})
}

// SetEx key to hold value in the cache and set key to timeout after the
// provided ttl.
func (c *Cache[K, V]) SetEx(key K, value V, ttl time.Duration) {
	expireAt := time.Now().Add(ttl)
	c.store.Set(key, Item[K, V]{Value: value, ExpireAt: &expireAt})
}

// Get returns the value associated with the provided key if it exists, or false
// if it does not.
//
// If the cache was opened with [WithPassiveExpiration] and the requested key
// is expired, it will be deleted from the cache and false will be returned.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	item, ok := c.store.Get(key, c.passiveExpiration)
	if !ok || item.IsExpired() {
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
	c.store.Delete(keys...)
}

// Flush the cache, deleting all keys.
func (c *Cache[K, V]) Flush() {
	c.store.Flush()
}

// Size returns the number of items currently in the cache.
func (c *Cache[K, V]) Size() int {
	return c.store.Size()
}

// Keys returns a map of all keys currently in the cache.
func (c *Cache[K, V]) Keys() []K {
	return c.store.Keys()
}

// Close the cache, stopping all running goroutines. Should be called when the
// cache is no longer needed.
func (c *Cache[K, V]) Close() {
	c.closer.Close()
}

func (c *Cache[K, V]) runActiveExpirer() {
	closerCh := c.closer.WaitClosed()
	ticker := time.NewTicker(c.activeExpirationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-closerCh:
			return
		case <-ticker.C:
			deleteAllExpiredKeys(c.store)
		}
	}
}

// deleteAllExpiredKeys from the provided store.
func deleteAllExpiredKeys[K comparable, V any](store storer[K, V]) {
	keys := store.Keys()
	expireKeys := make([]K, 0, len(keys))
	for _, key := range keys {
		if item, ok := store.Get(key, false); ok && item.IsExpired() {
			expireKeys = append(expireKeys, key)
		}
	}
	store.Delete(expireKeys...)
}

// storer is the interface depended upon by a cache.
type storer[K comparable, V any] interface {
	Set(key K, value Item[K, V])
	Get(key K, deleteExpired bool) (Item[K, V], bool)
	Delete(keys ...K)
	Items() (map[K]Item[K, V], unlockFunc)
	Keys() []K
	Size() int
	Flush()
}

// unlockFunc unlocks the mutex for the cache store.
type unlockFunc func()
