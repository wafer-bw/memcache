// Package memcache provides a generic in-memory key-value cache.
//
// The capacity of a cache is the total number of keys it is allowed to hold.
package memcache

import (
	"errors"
	"fmt"
	"time"

	"github.com/wafer-bw/memcache/internal/closeable"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/eviction/allkeyslfu"
	"github.com/wafer-bw/memcache/internal/eviction/allkeyslru"
	"github.com/wafer-bw/memcache/internal/eviction/noevict"
	"github.com/wafer-bw/memcache/internal/eviction/volatilelru"
	"github.com/wafer-bw/memcache/internal/expire"
	"github.com/wafer-bw/memcache/internal/ports"
)

var (
	ErrInvalidInterval = errors.New("provided interval must be greater than 0")
)

type InvalidCapacityError struct {
	Capacity int
	Minimum  int
	Policy   string
}

func (e InvalidCapacityError) Error() string {
	return fmt.Sprintf("capacity %d must be greater than %d for %s caches", e.Capacity, e.Minimum, e.Policy)
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
			return ErrInvalidInterval
		}
		c.expirer = expire.AllKeys[K, V]{}
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
		// because valid values for capacity depends on the policy, we do not
		// validate it here. Instead, it is done in each policy's open function.
		c.capacity = capacity
		return nil
	}
}

// Cache is a generic in-memory key-value cache.
type Cache[K comparable, V any] struct {
	closer                   ports.Closer
	store                    ports.Storer[K, V]
	expirer                  ports.Expirer[K, V]
	capacity                 int
	passiveExpiration        bool
	activeExpirationInterval time.Duration
}

// OpenNoEvictionCache opens a new in-memory key-value cache.
//
// This policy will ignore any additional keys that would cause the cache to
// breach its capacity.
//
// The capacity for this policy must be 0 (default) or set to a greater value
// via [WithCapacity].
func OpenNoEvictionCache[K comparable, V any](options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{
		closer:   closeable.New(),
		capacity: noevict.DefaultCapacity,
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if c.capacity < noevict.MinimumCapacity {
		return nil, InvalidCapacityError{
			Policy:   noevict.PolicyName,
			Capacity: c.capacity,
			Minimum:  noevict.MinimumCapacity,
		}
	}

	c.store = noevict.New[K, V](c.capacity)

	if c.activeExpirationInterval > 0 {
		go c.runActiveExpirer(c.activeExpirationInterval)
	}

	return c, nil
}

// OpenAllKeysLRUCache opens a new in-memory key-value cache.
//
// This policy evicts the least recently used key when the cache would breach
// its capacity.
//
// The capacity for this policy must be greater than 0.
func OpenAllKeysLRUCache[K comparable, V any](capacity int, options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{
		closer:   closeable.New(),
		capacity: capacity,
		expirer:  expire.AllKeys[K, V]{},
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if c.capacity < allkeyslru.MinimumCapacity {
		return nil, InvalidCapacityError{
			Policy:   allkeyslru.PolicyName,
			Capacity: c.capacity,
			Minimum:  allkeyslru.MinimumCapacity,
		}
	}

	c.store = allkeyslru.New[K, V](c.capacity)

	if c.activeExpirationInterval > 0 {
		go c.runActiveExpirer(c.activeExpirationInterval)
	}

	return c, nil
}

// OpenVolatileLRUCache opens a new in-memory key-value cache.
//
// This policy evicts the least recently used key with a ttl when the cache
// would breach its capacity. If no keys have a ttl, then the least recently
// used key is evicted.
//
// The capacity for this policy must be greater than 0.
func OpenVolatileLRUCache[K comparable, V any](capacity int, options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{
		closer:   closeable.New(),
		capacity: capacity,
		expirer:  expire.AllKeys[K, V]{},
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if c.capacity < volatilelru.MinimumCapacity {
		return nil, InvalidCapacityError{
			Policy:   volatilelru.PolicyName,
			Capacity: c.capacity,
			Minimum:  volatilelru.MinimumCapacity,
		}
	}

	c.store = volatilelru.New[K, V](c.capacity)

	if c.activeExpirationInterval > 0 {
		go c.runActiveExpirer(c.activeExpirationInterval)
	}

	return c, nil
}

// OpenAllKeysLFUCache opens a new in-memory key-value cache.
//
// This policy evicts the least frequently used key when the cache would breach
// its capacity.
//
// The capacity for this policy must be greater than 0.
func OpenAllKeysLFUCache[K comparable, V any](capacity int, options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{
		closer:   closeable.New(),
		capacity: capacity,
		expirer:  expire.AllKeys[K, V]{},
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if c.capacity < allkeyslfu.MinimumCapacity {
		return nil, InvalidCapacityError{
			Policy:   allkeyslfu.PolicyName,
			Capacity: c.capacity,
			Minimum:  allkeyslfu.MinimumCapacity,
		}
	}

	c.store = allkeyslfu.New[K, V](c.capacity)

	if c.activeExpirationInterval > 0 {
		go c.runActiveExpirer(c.activeExpirationInterval)
	}

	return c, nil
}

// Set non-expiring key to value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.store.Add(key, data.Item[K, V]{
		Value: value,
	})
}

// SetEx key that will expire after ttl to value in the cache.
func (c *Cache[K, V]) SetEx(key K, value V, ttl time.Duration) {
	expireAt := time.Now().Add(ttl)
	c.store.Add(key, data.Item[K, V]{
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
	if !ok {
		return *new(V), false
	}

	if item.IsExpired() {
		if c.passiveExpiration {
			c.store.Remove(key)
		}
		return *new(V), false
	}

	return item.Value, ok
}

// TTL for the provided key if it exists, or false if it does not. If the key is
// will not expire then (nil, true) will be returned.
func (c *Cache[K, V]) TTL(key K) (*time.Duration, bool) {
	item, ok := c.store.Get(key)
	return item.TTL(), ok
}

// Delete provided keys from the cache.
func (c *Cache[K, V]) Delete(keys ...K) {
	c.store.Remove(keys...)
}

// Size returns the number of items currently in the cache.
func (c *Cache[K, V]) Size() int {
	return c.store.Len()
}

// RandomKey returns a random key from the cache, or false if the cache is
// empty.
func (c *Cache[K, V]) RandomKey() (K, bool) {
	return c.store.RandomKey()
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
	c.closer.Close()
}

func (c *Cache[K, V]) closed() bool {
	return c.closer.Closed()
}

func (c *Cache[K, V]) runActiveExpirer(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.expirer.Expire(c)
		case <-c.closer.Ch():
			return
		}
	}
}
