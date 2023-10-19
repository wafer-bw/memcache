package memcache

// UnlockFunc unlockes the mutex for the cache store.
//
// export for testing.
type UnlockFunc func()

// export for testing.
func (c *Cache[K, V]) GetStore() (map[K]Item[K, V], UnlockFunc) {
	c.mu.Lock()

	return c.store, c.mu.Unlock
}

// export for testing.
func (c *Cache[K, V]) GetExpireOnGet() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.passiveExpiration
}

// export for testing.
func (c *Cache[K, V]) Closed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.closed
}
