package allkeyslru_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/store/allkeyslru"
)

var _ ports.Storer[int, int] = (*allkeyslru.Store[int, int])(nil)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a new store with provided capacity", func(t *testing.T) {
		t.Parallel()

		capacity := 10
		store := allkeyslru.New[int, int](capacity)
		require.Equal(t, capacity, store.Capacity())
	})

	t.Run("returns a new store with default capacity when provided an invalid one", func(t *testing.T) {
		t.Parallel()

		capacity := 1
		store := allkeyslru.New[int, int](capacity)
		require.Equal(t, allkeyslru.DefaultCapacity, store.Capacity())
	})
}

func TestStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value in all structures", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 10
		store := allkeyslru.New[int, int](2)
		store.Add(key, data.Item[int, int]{Value: val})

		elements, unlock := store.Elements()
		require.Len(t, elements, 1)
		require.Equal(t, key, elements[1].Value)
		unlock()

		list, unlock := store.List()
		require.Equal(t, 1, list.Len())
		require.Equal(t, key, list.Front().Value)
		unlock()

		items := store.Items()
		require.Len(t, items, 1)
		require.Equal(t, val, items[1].Value)
	})

	t.Run("evicts least recently used key when at capacity", func(t *testing.T) {
		t.Parallel()

		store := allkeyslru.New[int, int](2)
		store.Add(1, data.Item[int, int]{Value: 1})
		store.Add(2, data.Item[int, int]{Value: 2})
		_, _ = store.Get(1)
		store.Add(3, data.Item[int, int]{Value: 3})

		items := store.Items()
		require.Contains(t, items, 1)
		require.Contains(t, items, 3)
	})
}
