package data

import "time"

type Item[K comparable, V any] struct {
	Value    V
	ExpireAt *time.Time
	// TODO: Event methods (requires promoting package out of internal):
	//       They can cause a deadlock if they use the cache they are part of.
	//       - OnEvicted func(k K, v V)
	//       - OnExpired func(k K, v V)
	//       - OnDeleted func(k K, v V)
}

func (i Item[K, V]) IsExpired() bool {
	if i.ExpireAt == nil {
		return false
	}
	return time.Now().After(*i.ExpireAt)
}

func (i Item[K, V]) TTL() *time.Duration {
	if i.ExpireAt == nil {
		return nil
	}

	ttl := time.Until(*i.ExpireAt)
	if ttl < 0 {
		return new(time.Duration)
	}

	return &ttl
}
