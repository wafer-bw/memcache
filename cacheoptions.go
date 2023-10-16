package memcache

import "time"

type CacheOption[K comparable, V any] func(*Cache[K, V]) error

func WithPassiveExpiration[K comparable, V any]() CacheOption[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

// TODO: should interval be a separate option?
func WithExpirer[K comparable, V any](interval time.Duration, f ExpirerFunc[K, V]) CacheOption[K, V] {
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

// TODO: type EvictorFunc

// TODO: func WithEvictor
