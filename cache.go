package memcache

import (
	"time"

	"github.com/wafer-bw/memcache/errs"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/store/lru"
	"github.com/wafer-bw/memcache/internal/store/noevict"
)

// storer is the interface depended upon by a Cache.
type storer[K comparable, V any] interface {
	Set(key K, value data.Item[K, V])
	Get(key K) (data.Item[K, V], bool)
	Delete(keys ...K)
	Size() int
	Keys() []K
	Items() (map[K]data.Item[K, V], func())
	Flush()
	Close()
	Closed() bool
}

// Cache is a generic in-memory key-value cache.
type Cache[K comparable, V any] struct {
	store                    storer[K, V]
	capacity                 int
	passiveExpiration        bool
	activeExpirationInterval time.Duration
}

// OpenNoEvictionCache opens a new in-memory key-value cache using no eviction
// policy.
//
// This policy will ignore any additional keys that would cause the cache to
// breach its capacity.
//
// The capacity for this policy must be 0 (default) or set to a greater value
// via [WithCapacity].
func OpenNoEvictionCache[K comparable, V any](options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	var err error
	c.store, err = noevict.Open[K, V](noevict.Config{
		Capacity:                 c.capacity,
		PassiveExpiration:        c.passiveExpiration,
		ActiveExpirationInterval: c.activeExpirationInterval,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// OpenLRUCache opens a new in-memory key-value cache using a least recently
// used eviction policy.
//
// This policy evicts the least recently used key when the cache would breach
// its capacity.
//
// The capacity for this policy must be greater than 0.
func OpenLRUCache[K comparable, V any](capacity int, options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	var err error
	c.store, err = lru.Open[K, V](capacity, lru.Config{
		PassiveExpiration:        c.passiveExpiration,
		ActiveExpirationInterval: c.activeExpirationInterval,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Set non-expiring key to value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.store.Set(key, data.Item[K, V]{Value: value})
}

// SetEx key that will expire after ttl to value in the cache.
func (c *Cache[K, V]) SetEx(key K, value V, ttl time.Duration) {
	expireAt := time.Now().Add(ttl)
	c.store.Set(key, data.Item[K, V]{
		Value:    value,
		ExpireAt: &expireAt,
	})
}

// Get returns the value associated with the provided key if it exists, or false
// if it does not.
//
// If the cache was opened with [WithPassiveExpiration] and the requested key
// is expired, it will be deleted from the cache and false will be returned.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	item, ok := c.store.Get(key)
	return item.Value, ok
}

// Delete provided keys from the cache.
func (c *Cache[K, V]) Delete(keys ...K) {
	c.store.Delete(keys...)
}

// Size returns the number of items currently in the cache.
func (c *Cache[K, V]) Size() int {
	return c.store.Size()
}

// Keys returns a slice of all keys currently in the cache.
func (c *Cache[K, V]) Keys() []K {
	return c.store.Keys()
}

// Flush the cache, deleting all keys.
func (c *Cache[K, V]) Flush() {
	c.store.Flush()
}

// Close the cache, stopping all running goroutines. Should be called when the
// cache is no longer needed.
func (c *Cache[K, V]) Close() {
	c.store.Close()
}

// Option functions can be passed to open functions like [OpenNoEvictionCache]
// to control optional properties of the returned [Cache].
type Option[K comparable, V any] func(*Cache[K, V]) error

// WithPassiveExpiration enables the passive deletion of expired keys if they
// are found to be expired when accessed by [Cache.Get].
//
// This comes with a minor performance cost as [Cache.Get] now must acquire a
// write lock instead of a read lock.
func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

// WithActiveExpiration enables the active deletion of all expired keys at the
// provided interval.
//
// This comes with a minor performance cost as every tick requires a read lock
// to collect all expired keys followed by a write lock to delete them.
func WithActiveExpiration[K comparable, V any](interval time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) error {
		if interval <= 0 {
			return errs.ErrInvalidInterval
		}
		c.activeExpirationInterval = interval
		return nil
	}
}

// WithCapacity sets the maximum number of keys that the cache can hold.
//
// This option is made available to set the capacity of policies that do not
// need or use a capacity by default.
func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.capacity = int(capacity)
		return nil
	}
}
