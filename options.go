package memcache

import "time"

type Option[K comparable, V any] func(*Cache[K, V]) error

func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

func WithExpirer[K comparable, V any](f ExpirerFunc[K, V], interval time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) error {
		if f == nil {
			return nil // TODO: return error
		} else if interval <= 0 {
			return nil // TODO: return error
		}

		c.expirationInterval = interval
		c.expirer = f
		return nil
	}
}

// TODO: func WithDefaultExpirer

// TODO: type EvictorFunc

// TODO: func WithEvictor

// TODO: func WithDefaultEvictor
