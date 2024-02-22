package expire

import (
	"time"
)

const (
	DefaultSampleSize    int     = 20
	DefaultExpirePercent float32 = 0.25
)

// Cacher is the interface depended upon by an expirer.
type Cacher[K comparable, V any] interface {
	TTL(key K) (*time.Duration, bool)
	Delete(keys ...K)
	Keys() []K
	Size() int
	RandomKey() (K, bool)
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

type RandomSample[K comparable, V any] struct {
	SampleSize    int
	ExpirePercent float32
}

func (e *RandomSample[K, V]) Expire(cache Cacher[K, V]) {
	if e.SampleSize <= 0 {
		e.SampleSize = DefaultSampleSize
	}

	if e.ExpirePercent <= 0 {
		e.ExpirePercent = DefaultExpirePercent
	}

	s := cache.Size()
	if s == 0 {
		return
	}

	expiredCount := 0
	for i := 0; i < e.SampleSize; i++ {
		key, ok := cache.RandomKey()
		if !ok {
			return
		}

		if ttl, ok := cache.TTL(key); ok && ttl != nil && *ttl <= 0 {
			expiredCount++
			cache.Delete(key)
		}
	}

	percentExpired := float32(expiredCount) / float32(e.SampleSize)

	if percentExpired > e.ExpirePercent {
		e.Expire(cache)
	}
}
