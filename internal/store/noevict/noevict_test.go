package noevict_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/store/noevict"
)

var _ ports.Storer[int, int] = (*noevict.Store[int, int])(nil)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a new store with provided capacity", func(t *testing.T) {
		t.Parallel()

		capacity := 10
		store := noevict.New[int, int](capacity)
		require.Equal(t, capacity, store.Capacity())
	})

	t.Run("returns a new store with default capacity when provided an invalid one", func(t *testing.T) {
		t.Parallel()

		capacity := -1
		store := noevict.New[int, int](capacity)
		require.Equal(t, noevict.DefaultCapacity, store.Capacity())
	})
}

func TestStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 10
		store := noevict.New[int, int](0)
		store.Add(key, data.Item[int, int]{Value: val})

		items := store.Items()
		require.Len(t, items, 1)
		require.Equal(t, val, items[key].Value)
	})

	t.Run("does not add more keys when at capacity", func(t *testing.T) {
		t.Parallel()

		store := noevict.New[int, int](2)
		store.Add(1, data.Item[int, int]{Value: 10})
		store.Add(1, data.Item[int, int]{Value: 10})
		store.Add(1, data.Item[int, int]{Value: 10})
		store.Add(2, data.Item[int, int]{Value: 20})
		store.Add(3, data.Item[int, int]{Value: 30})

		items := store.Items()
		require.Len(t, items, 2)
		require.Contains(t, items, 1)
		require.Contains(t, items, 2)
	})
}
