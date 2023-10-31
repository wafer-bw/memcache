package memcache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

func TestLRUStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value in all structures", func(t *testing.T) {
		t.Parallel()

		store, _ := memcache.NewLRUStore[int, int](2)
		store.Underlying.Set(1, memcache.Item[int, int]{Value: 1})

		require.Len(t, store.Items(), 1)
		require.Len(t, store.Elements(), 1)
		require.Equal(t, 1, store.List().Len())
		require.Len(t, store.Keys(), 1)

		require.Equal(t, 1, store.Items()[1].Value)
		require.Equal(t, 1, store.Elements()[1].Value)
		require.Equal(t, store.Elements()[1], store.List().Front())
		require.Equal(t, struct{}{}, store.Keys()[1])
	})

	t.Run("evicts least recently used key from all structures", func(t *testing.T) {
		t.Parallel()

		store, _ := memcache.NewLRUStore[int, int](2)
		store.Underlying.Set(1, memcache.Item[int, int]{Value: 1})
		store.Underlying.Set(2, memcache.Item[int, int]{Value: 2})
		_, _ = store.Underlying.Get(1, false)
		store.Underlying.Set(3, memcache.Item[int, int]{Value: 3})

		require.Len(t, store.Items(), 2)
		require.Len(t, store.Elements(), 2)
		require.Equal(t, 2, store.List().Len())
		require.Len(t, store.Keys(), 2)

		require.Equal(t, 1, store.Items()[1].Value)
		require.Equal(t, 3, store.Items()[3].Value)
		require.Equal(t, 1, store.Elements()[1].Value)
		require.Equal(t, 3, store.Elements()[3].Value)
		require.Equal(t, store.Elements()[1], store.List().Front().Next())
		require.Equal(t, store.Elements()[3], store.List().Front())
		require.Equal(t, struct{}{}, store.Keys()[1])
		require.Equal(t, struct{}{}, store.Keys()[3])
	})

	t.Run("returns an error if capacity is lower than 2", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.NewLRUStore[int, int](1)
		require.ErrorIs(t, err, memcache.ErrInvalidCapacity)
	})
}
