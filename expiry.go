package memcache

type expirer[K comparable, V any] interface {
	Expire(cache *Cache[K, V])
}

type expirerFunc[K comparable, V any] func(cache *Cache[K, V])

func (f expirerFunc[K, V]) Expire(cache *Cache[K, V]) {
	f(cache)
}

func fullScanExpirer[K comparable, V any](cache *Cache[K, V]) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for key, record := range cache.store {
		if record.IsExpired() {
			delete(cache.store, key)
		}
	}
}
