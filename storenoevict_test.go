package memcache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

func TestNoEvictStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value in all structures", func(t *testing.T) {
		t.Parallel()

		store := memcache.NewNoEvictStore[int, int]()
		store.Underlying.Set(1, memcache.Item[int, int]{Value: 1})

		require.Len(t, store.Items(), 1)
		require.Len(t, store.Keys(), 1)

		require.Equal(t, 1, store.Items()[1].Value)
		require.Equal(t, struct{}{}, store.Keys()[1])
	})
}
