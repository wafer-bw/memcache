package memcache

import "time"

type Item[K comparable, V any] struct {
	Value    V
	ExpireAt *time.Time
	// OnEvicted func(k K, v V) // TODO: add this.
}

func (i Item[K, V]) IsExpired() bool {
	if i.ExpireAt == nil {
		return false
	}

	return time.Now().After(*i.ExpireAt)
}
