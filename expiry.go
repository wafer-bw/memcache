package memcache

type fullScanExpirer[K comparable, V any] struct{}

func (f *fullScanExpirer[K, V]) Expire(cache *Cache[K, V]) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for key, record := range cache.store {
		if record.IsExpired() {
			delete(cache.store, key)
		}
	}
}
