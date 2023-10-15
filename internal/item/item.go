package item

import "time"

type Item[V any] struct {
	Value    V
	ExpireAt *time.Time
}

func (r Item[V]) IsExpired() bool {
	if r.ExpireAt == nil {
		return false
	}

	return time.Now().After(*r.ExpireAt)
}
