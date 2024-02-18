package expire

import (
	"time"
)

type Storer[K comparable, V any] interface {
	TTL(key K) (*time.Duration, bool)
	Delete(keys ...K)
	Keys() []K
}

type AllKeys[K comparable, V any] struct {
}

func (e AllKeys[K, V]) Expire(store Storer[K, V]) {
	keys := store.Keys()
	for _, key := range keys {
		if ttl, ok := store.TTL(key); ok && ttl != nil && *ttl <= 0 {
			store.Delete(key)
		}
	}
}
