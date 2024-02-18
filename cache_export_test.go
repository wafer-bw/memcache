package memcache

import (
	"time"
)

// export for testing.
type Storer[K comparable, V any] storer[K, V]

// export for testing.
func (c *Cache[K, V]) Store() Storer[K, V] {
	return c.store
}

// export for testing.
func (c *Cache[K, V]) PassiveExpiration() bool {
	return c.passiveExpiration
}

// export for testing.
func (c *Cache[K, V]) ExpirationInterval() time.Duration {
	return c.activeExpirationInterval
}

// export for testing.
func (c *Cache[K, V]) Closed() bool {
	return c.store.Closed()
}
