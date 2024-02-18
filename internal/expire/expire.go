package expire

import "github.com/wafer-bw/memcache/internal/data"

type Storer[K comparable, V any] interface {
	Get(key K) (data.Item[K, V], bool)
	Delete(keys ...K)
	Keys() []K
}

type AllKeys[K comparable, V any] struct {
}

func (e AllKeys[K, V]) Expire(store Storer[K, V]) {
	keys := store.Keys()
	for _, key := range keys {
		if item, ok := store.Get(key); ok && item.IsExpired() {
			store.Delete(key)
		}
	}
}
