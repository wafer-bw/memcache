package expire_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/expire"
)

type expirer[K comparable, V any] interface {
	Expire(expire.Storer[K, V])
}

var _ expirer[int, int] = (*expire.AllKeys[int, int])(nil)

type store[K comparable, V any] struct {
	items map[K]data.Item[K, V]
}

func (s *store[K, V]) Get(key K) (data.Item[K, V], bool) {
	item, ok := s.items[key]
	return item, ok
}

func (s *store[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(s.items, key)
	}
}

func (s *store[K, V]) Keys() []K {
	keys := make([]K, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	return keys
}

func TestAllKeys_Expire(t *testing.T) {
	t.Parallel()

	t.Run("deletes all expired items", func(t *testing.T) {
		t.Parallel()

		expired := time.Now().Add(-1 * time.Minute)
		unexpired := time.Now().Add(1 * time.Minute)

		s := &store[int, int]{
			items: map[int]data.Item[int, int]{
				1: {ExpireAt: nil},
				2: {ExpireAt: &expired},
				3: {ExpireAt: &expired},
				4: {ExpireAt: &unexpired},
				5: {ExpireAt: nil},
				6: {ExpireAt: &unexpired},
				7: {ExpireAt: &expired},
			},
		}

		e := expire.AllKeys[int, int]{}
		e.Expire(s)
		require.Equal(t, 4, len(s.items))
	})
}
