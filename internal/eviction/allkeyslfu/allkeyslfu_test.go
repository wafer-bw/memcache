package allkeyslfu_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/eviction/allkeyslfu"
	"github.com/wafer-bw/memcache/internal/ports"
)

var _ ports.Storer[int, int] = (*allkeyslfu.Store[int, int])(nil)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a new store with provided capacity", func(t *testing.T) {
		t.Parallel()

		capacity := 10
		store := allkeyslfu.New[int, int](capacity)
		require.Equal(t, capacity, store.Capacity())
	})

	t.Run("returns a new store with default capacity when provided an invalid one", func(t *testing.T) {
		t.Parallel()

		capacity := 1
		store := allkeyslfu.New[int, int](capacity)
		require.Equal(t, allkeyslfu.DefaultCapacity, store.Capacity())
	})
}

func TestStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 10
		store := allkeyslfu.New[int, int](0)
		store.Add(key, data.Item[int, int]{Value: val})

		items := store.Items()
		require.Len(t, items, 1)
		require.Equal(t, val, items[key].Value)
	})

	t.Run("evicts least frequently used key when at capacity", func(t *testing.T) {
		t.Parallel()

		store := allkeyslfu.New[int, int](3)
		store.Add(1, data.Item[int, int]{Value: 1})
		store.Add(2, data.Item[int, int]{Value: 2})
		store.Add(3, data.Item[int, int]{Value: 3})
		_, _ = store.Get(1)
		_, _ = store.Get(1)
		_, _ = store.Get(2)
		_, _ = store.Get(2)
		_, _ = store.Get(3)

		items := store.Items()
		require.Contains(t, items, 1)
		require.Contains(t, items, 2)
	})

	t.Run("does not add more keys when at capacity", func(t *testing.T) {
		t.Parallel()

		store := allkeyslfu.New[int, int](2)
		store.Add(1, data.Item[int, int]{Value: 10})
		store.Add(1, data.Item[int, int]{Value: 10})
		store.Add(1, data.Item[int, int]{Value: 10})
		store.Add(2, data.Item[int, int]{Value: 20})
		store.Add(3, data.Item[int, int]{Value: 30})

		items := store.Items()
		require.Len(t, items, 2)
	})
}

func TestStore_Flush(t *testing.T) {
	t.Parallel()

	t.Run("clears all keys and values", func(t *testing.T) {
		t.Parallel()

		store := allkeyslfu.New[int, int](2)
		store.Add(1, data.Item[int, int]{Value: 1})
		store.Add(2, data.Item[int, int]{Value: 2})
		store.Flush()

		items := store.Items()
		require.Empty(t, items)
	})
}
