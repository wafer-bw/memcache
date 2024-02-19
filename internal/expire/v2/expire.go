package expire

import "time"

type Cacher[K comparable, V any] interface {
	TTL(key K) (*time.Duration, bool)
	Delete(keys ...K)
	Keys() []K
}

type AllKeys[K comparable, V any] struct {
}

func (e AllKeys[K, V]) Expire(cache Cacher[K, V]) {
	keys := cache.Keys()
	for _, key := range keys {
		if ttl, ok := cache.TTL(key); ok && ttl != nil && *ttl <= 0 {
			cache.Delete(key)
		}
	}
}
