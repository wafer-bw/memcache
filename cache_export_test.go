package memcache

// Everything in this file is exported for testing purposes only.

import (
	"time"
)

func (c *Cache[K, V]) Store() storer[K, V] {
	return c.store
}

func (c *Cache[K, V]) PassiveExpiration() bool {
	return c.passiveExpiration
}

func (c *Cache[K, V]) ExpirationInterval() time.Duration {
	return c.activeExpirationInterval
}

func (c *Cache[K, V]) Closed() bool {
	select {
	case <-c.closeCh:
		return true
	default:
		return false
	}
}

func DeleteAllExpiredKeys[K comparable, V any](store storer[K, V]) {
	deleteAllExpiredKeys[K, V](store)
}
