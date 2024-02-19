package memcache

import (
	"time"

	"github.com/wafer-bw/memcache/internal/ports"
)

// export for testing.
func (c *Cache[K, V]) Store() ports.Storer[K, V] {
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
//
// TODO: remove.
func (c *Cache[K, V]) Capacity() int {
	return c.capacity
}

// export for testing.
func (c *Cache[K, V]) Closed() bool {
	return c.closed()
}
