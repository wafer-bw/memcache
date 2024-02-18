package memcache

import (
	"time"

	"github.com/wafer-bw/memcache/errs"
)

// Option functions can be passed to [Open] to control optional properties of
// the returned [Cache].
type Option[K comparable, V any] func(*Cache[K, V]) error

// WithPassiveExpiration enables the passive deletion of expired keys if they
// are found to be expired when accessed by [Cache.Get].
func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

// WithActiveExpiration enables the active deletion of expired keys at the
// provided interval.
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
// This option is made available to control the capacity of policies that do not
// require a capacity.
func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.capacity = capacity
		return nil
	}
}
