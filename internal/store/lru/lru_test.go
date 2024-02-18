package lru_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/store/lru"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	t.Run("returns a new store with minimum configuration", func(t *testing.T) {
		t.Parallel()

		store, err := lru.Open[int, int](1, lru.Config{})
		require.NoError(t, err)
		require.NotNil(t, store)
	})

	t.Run("returns an error if capacity is lower than  the minimum", func(t *testing.T) {
		t.Parallel()

		_, err := lru.Open[int, int](lru.MinimumCapacity-1, lru.Config{})
		require.Error(t, err)
	})
}

func TestStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 100
		store, _ := lru.Open[int, int](1, lru.Config{})
		store.Set(key, data.Item[int, int]{Value: val})

		items, unlock := store.Items()
		defer unlock()
		require.Len(t, items, 1)
		require.Equal(t, val, items[key].Value)
	})

	t.Run("stores key and value in all structures", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 10
		store, _ := lru.Open[int, int](2, lru.Config{})
		store.Set(key, data.Item[int, int]{Value: val})

		items, unlock := store.Items()
		require.Len(t, items, 1)
		require.Equal(t, val, items[1].Value)
		unlock()

		elements, unlock := store.Elements()
		require.Len(t, elements, 1)
		require.Equal(t, key, elements[1].Value)
		unlock()

		list, unlock := store.List()
		require.Equal(t, 1, list.Len())
		require.Equal(t, key, list.Front().Value)
		unlock()
	})

	t.Run("evicts least recently used key from all structures", func(t *testing.T) {
		t.Parallel()

		store, _ := lru.Open[int, int](2, lru.Config{})
		store.Set(1, data.Item[int, int]{Value: 1})
		store.Set(2, data.Item[int, int]{Value: 2})
		_, _ = store.Get(1)
		store.Set(3, data.Item[int, int]{Value: 3})

		items, unlock := store.Items()
		require.Len(t, items, 2)
		require.Contains(t, items, 1)
		require.Contains(t, items, 3)
		unlock()
	})
}
