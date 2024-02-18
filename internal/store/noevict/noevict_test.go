package noevict_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/store/noevict"
)

func TestStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 10
		store, _ := noevict.Open[int, int](noevict.Config{})
		store.Set(key, data.Item[int, int]{Value: val})

		items, unlock := store.Items()
		defer unlock()
		require.Len(t, items, 1)
		require.Equal(t, val, items[key].Value)
	})

	t.Run("does not add more keys when at capacity", func(t *testing.T) {
		t.Parallel()

		store, _ := noevict.Open[int, int](noevict.Config{Capacity: 2})
		store.Set(1, data.Item[int, int]{Value: 10})
		store.Set(1, data.Item[int, int]{Value: 10})
		store.Set(1, data.Item[int, int]{Value: 10})
		store.Set(2, data.Item[int, int]{Value: 20})
		store.Set(3, data.Item[int, int]{Value: 30})

		items, unlock := store.Items()
		defer unlock()
		require.Len(t, items, 2)
		require.Contains(t, items, 1)
		require.Contains(t, items, 2)
	})

	t.Run("returns an error if capacity is lower than the minimum", func(t *testing.T) {
		t.Parallel()

		_, err := noevict.Open[int, int](noevict.Config{Capacity: noevict.MinimumCapacity - 1})
		require.Error(t, err)
	})
}
