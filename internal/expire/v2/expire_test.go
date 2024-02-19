package expire_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/expire/v2"
	"github.com/wafer-bw/memcache/internal/ports"
)

var _ ports.Expirer[int, int] = (*expire.AllKeys[int, int])(nil)

type cache[K comparable, V any] struct {
	items map[K]data.Item[K, V]
}

func (c *cache[K, V]) TTL(key K) (*time.Duration, bool) {
	item, ok := c.items[key]
	return item.TTL(), ok
}

func (c *cache[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(c.items, key)
	}
}

func (c *cache[K, V]) Keys() []K {
	keys := make([]K, 0, len(c.items))
	for key := range c.items {
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

		c := &cache[int, int]{
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
		e.Expire(c)
		require.Equal(t, 4, len(c.items))
	})
}
