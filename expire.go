package memcache

type ExpirerFunc[K comparable, V any] func(store map[K]Item[K, V])

func DeleteAllExpiredKeys[K comparable, V any](store map[K]Item[K, V]) {
	for k, v := range store {
		if v.IsExpired() {
			delete(store, k)
		}
	}
}
