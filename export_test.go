package memcache

import (
	"github.com/wafer-bw/memcache/internal/item"
)

// UnlockFunc unlockes the mutex for the cache store.
//
// export for testing.
type UnlockFunc func()

// export for testing.
func (c *Cache[K, V]) GetStore() (map[K]item.Item[V], UnlockFunc) {
	c.mu.Lock()

	return c.store, c.mu.Unlock
}

// export for testing.
func (c *Cache[K, V]) GetExpireOnGet() bool {
	return c.passiveExpiration
}
