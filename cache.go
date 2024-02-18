package memcache

import (
	"time"

	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/store/lru"
	"github.com/wafer-bw/memcache/internal/store/noevict"
)

// storer is the interface depended upon by a cache.
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

type policy string

const (
	policyNoEvict policy = policy(noevict.PolicyName)
	policyLRU     policy = policy(lru.PolicyName)
)

// Cache is a generic in-memory key-value cache.
type Cache[K comparable, V any] struct {
	store                    storer[K, V]
	capacity                 int
	policy                   policy
	passiveExpiration        bool
	activeExpirationInterval time.Duration
}

// OpenNoEvictionCache opens a new in-memory key-value cache using no eviction
// policy.
//
// This policy will ignore any additional keys that would cause the cache to
// breach its capacity.
//
// The capacity for this policy must be 0 (default) or greater via
// [WithCapacity].
func OpenNoEvictionCache[K comparable, V any](options ...Option[K, V]) (*Cache[K, V], error) {
	c := &Cache[K, V]{policy: policyNoEvict}
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
	c := &Cache[K, V]{policy: policyLRU}
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

// Set permanent key to hold value in the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.store.Set(key, data.Item[K, V]{Value: value})
}

// SetEx key to hold value in the cache and set key to timeout after the
// provided ttl.
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

// Keys returns a map of all keys currently in the cache.
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
