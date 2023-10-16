package memcache

import "time"

type ItemOption[K comparable, V any] func(*Item[K, V])

func WithTTL[K comparable, V any](d time.Duration) ItemOption[K, V] {
	return func(i *Item[K, V]) {
		expireAt := time.Now().Add(d)
		i.ExpireAt = &expireAt
	}
}

// TODO: WithOnEvicted(func(k K, v V) { ... }
