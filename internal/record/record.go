package record

import "time"

type Record[V any] struct {
	Value    V
	ExpireAt *time.Time
}

func (r Record[V]) Expired() bool {
	if r.ExpireAt == nil {
		return false
	}

	return time.Now().After(*r.ExpireAt)
}
