package lru_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/errs"
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

	t.Run("returns an error if capacity is lower than 1", func(t *testing.T) {
		t.Parallel()

		_, err := lru.Open[int, int](lru.MinimumCapacity-1, lru.Config{})
		require.ErrorAs(t, err, &errs.InvalidCapacityError{})
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

	// t.Run("evicts least recently used key from all structures", func(t *testing.T) {
	// 	t.Parallel()

	// 	store, _ := lru.Open[int, int](2)
	// 	store.Underlying.Set(1, data.Item[int, int]{Value: 1})
	// 	store.Underlying.Set(2, data.Item[int, int]{Value: 2})
	// 	_, _ = store.Underlying.Get(1, false)
	// 	store.Underlying.Set(3, data.Item[int, int]{Value: 3})

	// 	require.Len(t, store.Items(), 2)
	// 	require.Len(t, store.Elements(), 2)
	// 	require.Equal(t, 2, store.List().Len())

	// 	require.Equal(t, 1, store.Items()[1].Value)
	// 	require.Equal(t, 3, store.Items()[3].Value)
	// 	require.Equal(t, 1, store.Elements()[1].Value)
	// 	require.Equal(t, 3, store.Elements()[3].Value)
	// 	require.Equal(t, store.Elements()[1], store.List().Front().Next())
	// 	require.Equal(t, store.Elements()[3], store.List().Front())
	// })

}
