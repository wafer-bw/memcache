package memcache

import (
	"time"

	"github.com/wafer-bw/memcache/errs"
	"github.com/wafer-bw/memcache/internal/store/lru"
)

// Option defines the signature of a function that can be passed to [Open] as
// a functional option for controlling the behavior of the returned [Cache]
type Option[K comparable, V any] func(*Cache[K, V]) error

// WithPassiveExpiration enables the passive deletion of expired keys via
// read methods such as [Cache.Get].
//
// This can be combined with [WithActiveExpiration] to enable both passive and
// active expiration of keys. See [Open] for example usage.
func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

// WithActiveExpiration enables the active deletion of expired keys.
//
// This can be combined with [WithPassiveExpiration] to enable both passive and
// active expiration of keys.
//
// See [Open] for example usage.
func WithActiveExpiration[K comparable, V any](interval time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) error {
		if interval <= 0 {
			return errs.ErrInvalidInterval
		}
		c.activeExpirationInterval = interval
		return nil
	}
}

// WithLRUEviction enables the eviction of the least recently used key when the
// cache would breach its capacity.
//
// Calculating the size of a generic map in memory incurrs a heavy performance
// cost. For that reason, the capacity of a cache is defined as the total number
// of keys it is allowed to hold.
func WithLRUEviction[K comparable, V any](capacity int) Option[K, V] {
	return func(c *Cache[K, V]) error {
		if capacity < lru.MinimumCapacity {
			return errs.InvalidCapacityError{
				Policy:   lru.PolicyName,
				Capacity: capacity,
				Minimum:  lru.MinimumCapacity,
			}
		}

		c.policy = policyLRU
		c.capacity = capacity
		return nil
	}
}

// WithCapacity is used to set the capacity of the cache if the chosen or
// default policy does not otherwise require one.
func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(c *Cache[K, V]) error {
		if capacity < 0 {
			return errs.InvalidCapacityError{
				Policy:   "all",
				Capacity: capacity,
				Minimum:  0,
			}
		}
		c.capacity = capacity
		return nil
	}
}
