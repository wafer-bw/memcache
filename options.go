package memcache

import (
	"time"
)

type Option[K comparable, V any] func(*Cache[K, V]) error

func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

func WithDefaultExpirer[K comparable, V any](interval time.Duration) Option[K, V] {
	return WithExpirer[K, V](DeleteAllExpiredKeys, interval)
}

func WithExpirer[K comparable, V any](f ExpirerFunc[K, V], interval time.Duration) Option[K, V] {
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
