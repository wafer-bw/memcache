package volatilelru_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/eviction/volatilelru"
	"github.com/wafer-bw/memcache/internal/ports"
)

var _ ports.Storer[int, int] = (*volatilelru.Store[int, int])(nil)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a new store with provided capacity", func(t *testing.T) {
		t.Parallel()

		capacity := 10
		store := volatilelru.New[int, int](capacity)
		require.Equal(t, capacity, store.Capacity())
	})

	t.Run("returns a new store with default capacity when provided an invalid one", func(t *testing.T) {
		t.Parallel()

		capacity := 1
		store := volatilelru.New[int, int](capacity)
		require.Equal(t, volatilelru.DefaultCapacity, store.Capacity())
	})
}

func TestStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("stores key and value in all structures", func(t *testing.T) {
		t.Parallel()

		key, val := 1, 10
		store := volatilelru.New[int, int](2)
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

	t.Run("evicts least recently used key with a ttl when at capacity", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now().Add(1 * time.Hour)
		store := volatilelru.New[int, int](3)
		store.Add(1, data.Item[int, int]{Value: 1})
		store.Add(2, data.Item[int, int]{Value: 2, ExpireAt: &expireAt})
		store.Add(3, data.Item[int, int]{Value: 3, ExpireAt: &expireAt})
		_, _ = store.Get(2)
		store.Add(4, data.Item[int, int]{Value: 4})

		items := store.Items()
		require.Contains(t, items, 1)
		require.Contains(t, items, 2)
		require.Contains(t, items, 4)
	})

	t.Run("evicts least recently used key when at capacity no keys have a ttl", func(t *testing.T) {
		t.Parallel()

		store := volatilelru.New[int, int](2)
		store.Add(1, data.Item[int, int]{Value: 1})
		store.Add(2, data.Item[int, int]{Value: 2})
		_, _ = store.Get(1)
		store.Add(3, data.Item[int, int]{Value: 3})

		items := store.Items()
		require.Contains(t, items, 1)
		require.Contains(t, items, 3)
	})
}

func TestStore_Flush(t *testing.T) {
	t.Parallel()

	t.Run("clears all keys and values", func(t *testing.T) {
		t.Parallel()

		store := volatilelru.New[int, int](2)
		store.Add(1, data.Item[int, int]{Value: 1})
		store.Add(2, data.Item[int, int]{Value: 2})
		store.Flush()

		items := store.Items()
		require.Empty(t, items)

		elements, unlock := store.Elements()
		require.Empty(t, elements)
		unlock()

		list, unlock := store.List()
		require.Equal(t, 0, list.Len())
		unlock()
	})
}
