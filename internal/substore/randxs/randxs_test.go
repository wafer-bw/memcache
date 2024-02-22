package randxs_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/substore/randxs"
)

var _ ports.RandomAccessor[int] = (*randxs.Store[int])(nil)

func TestStore_Add(t *testing.T) {
	t.Parallel()

	t.Run("adds key to store", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)

		keys, unlock := store.Keys()
		require.Equal(t, []int{1}, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Equal(t, map[int]int{1: 0}, keyIndices)
		unlock()
	})

	t.Run("grows beyond capacity", func(t *testing.T) {
		t.Parallel()

		capacity := 2
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(2)
		store.Add(3)

		keys, unlock := store.Keys()
		require.Equal(t, []int{1, 2, 3}, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Equal(t, map[int]int{1: 0, 2: 1, 3: 2}, keyIndices)
		unlock()
	})

	t.Run("does not corrupt store on subsequent identical adds", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(1)

		keys, unlock := store.Keys()
		require.Equal(t, []int{1}, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Equal(t, map[int]int{1: 0}, keyIndices)
		unlock()
	})
}

func TestStore_RandomKey(t *testing.T) {
	t.Parallel()

	t.Run("returns random key", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(2)
		store.Add(3)
		store.Add(4)

		keysFound := make(map[int]struct{}, capacity)
		for i := 0; i < 100; i++ {
			key, ok := store.RandomKey()
			require.True(t, ok)
			keysFound[key] = struct{}{}
		}
		require.Len(t, keysFound, capacity)
		require.Equal(t, keysFound, map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}})
	})

	t.Run("returns zero value and false when empty", func(t *testing.T) {
		t.Parallel()

		store := randxs.New[int](1)

		key, ok := store.RandomKey()
		require.Zero(t, key)
		require.False(t, ok)
	})
}

func TestStore_Remove(t *testing.T) {
	t.Parallel()

	t.Run("removes key", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(2)

		store.Remove(1)
		keys, unlock := store.Keys()
		require.Equal(t, []int{2}, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Equal(t, map[int]int{2: 0}, keyIndices)
		unlock()
	})

	t.Run("does not corrupt store on subsequent identical removes", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(2)

		store.Remove(1)
		store.Remove(1)
		keys, unlock := store.Keys()
		require.Equal(t, []int{2}, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Equal(t, map[int]int{2: 0}, keyIndices)
		unlock()
	})

	t.Run("correctly adjusts key indices after removal", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(2)
		store.Add(3)
		store.Add(4)

		store.Remove(2)
		keys, unlock := store.Keys()
		require.Equal(t, []int{1, 4, 3}, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Equal(t, map[int]int{1: 0, 3: 2, 4: 1}, keyIndices)
		unlock()
	})
}

func TestStore_Clear(t *testing.T) {
	t.Parallel()

	t.Run("removes all keys from all structures", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := randxs.New[int](capacity)
		store.Add(1)
		store.Add(2)
		store.Add(3)
		store.Add(4)

		store.Clear()
		keys, unlock := store.Keys()
		require.Empty(t, keys)
		unlock()

		keyIndices, unlock := store.KeyIndices()
		require.Empty(t, keyIndices)
		unlock()
	})
}
