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
// active expiration of keys.
//
// See [Open] for example usage.
func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

// WithActiveExpiration enables the active deletion of expired keys by running
// the provided expirer function at the provided interval.
//
// This can be combined with [WithPassiveExpiration] to enable both passive and
// active expiration of keys.
//
// See [Open] for example usage.
func WithActiveExpiration[K comparable, V any](f ExpirerFunc[K, V], interval time.Duration) Option[K, V] {
	// TODO: consider accepting an interface instead of a type

	return func(c *Cache[K, V]) error {
		if f == nil {
			return ErrNilExpirerFunc
		} else if interval <= 0 {
			return ErrInvalidInterval
		}

		c.expirationInterval = interval
		c.expirer = f
		return nil
	}
}

// TODO: type EvictorFunc

// TODO: func WithEvictor
