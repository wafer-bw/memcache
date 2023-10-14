package memcache

import (
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
func (c *CacheConfig) GetEvictionInterval() time.Duration {
	return c.evictionInterval
}

// export for testing.
func (c *CacheConfig) GetExpirationInterval() time.Duration {
	return c.expirationInterval
}

func (c *ValueConfig) GetExpireAt() *time.Time {
	return c.expireAt
}
