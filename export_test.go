package memcache

import (
	"github.com/wafer-bw/memcache/internal/record"
)

// UnlockFunc unlockes the mutex for the cache store.
//
// export for testing.
type UnlockFunc func()

// export for testing.
func (c *Cache[K, V]) GetStore() (map[K]record.Record[V], UnlockFunc) {
	c.mu.Lock()

	return c.store, c.mu.Unlock
}

// export for testing.
func (c *Cache[K, V]) GetExpireOnGet() bool {
	return c.passiveExpiration
}
