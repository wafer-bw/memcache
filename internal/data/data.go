package data

import "time"

type Item[K comparable, V any] struct {
	Value    V
	ExpireAt *time.Time
	// TODO: OnEvicted func(k K, v V)
	// TODO: OnExpired func(k K, v V)
	// TODO: OnDeleted func(k K, v V)
}

func (i Item[K, V]) IsExpired() bool {
	if i.ExpireAt == nil {
		return false
	}
	return time.Now().After(*i.ExpireAt)
}
