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
func (c *Cache[K, V]) Lock() {
	c.mu.Lock()
}

// export for testing.
func (c *Cache[K, V]) Unlock() {
	c.mu.Unlock()
}

// export for testing.
func (c *Cache[K, V]) RLock() {
	c.mu.RLock()
}

// export for testing.
func (c *Cache[K, V]) RUnlock() {
	c.mu.RUnlock()
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

// export for testing.
func (c *Cache[K, V]) GetExpirer() ExpirerFunc[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.expirer
}
