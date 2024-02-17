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
	return c.store.Closed()
}
