package memcache

import (
	"time"
)

// Option defines the signature of a function that can be passed to [Open] as
// a functional option for controlling the behavior of the returned [Cache]
type Option[K comparable, V any] func(*Cache[K, V]) error

// WithPassiveExpiration enables the passive deletion of expired keys via
// read methods such as [Cache.Get] & [Cache.Has].
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
			return ErrInvalidInterval
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
		store, err := newLRUStore[K, V](capacity)
		if err != nil {
			return err
		}

		c.store = store

		return nil
	}
}
