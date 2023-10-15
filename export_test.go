package memcache

import (
	"context"
	"sync"
	"time"

	"github.com/wafer-bw/memcache/internal/record"
)

// export for testing.
func (c *Cache[K, V]) GetStore() map[K]record.Record[V] {
	return c.store
}

// export for testing.
func (c *Cache[K, V]) GetMutex() *sync.RWMutex {
	return &c.mu
}

// export for testing.
func (c *Cache[K, V]) GetExpirationInterval() time.Duration {
	return c.expirationInterval
}

// export for testing.
func (c *Cache[K, V]) GetExpireOnGet() bool {
	return c.expireOnGet
}

// export for testing.
func (c *Cache[K, V]) RunExpirer(ctx context.Context) {
	c.runExpirer(ctx)
}
