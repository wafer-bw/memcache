package lfulist_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/substore/lfulist"
)

var _ ports.LFUTracker[int] = (*lfulist.Store[int])(nil)

func TestStore(t *testing.T) {
	t.Parallel()

	t.Run("returns least frequently used key", func(t *testing.T) {
		t.Parallel()

		capacity := 4
		store := lfulist.New[int](capacity)
		store.Inc(1)
		store.Inc(2)
		store.Inc(3)
		store.Inc(1)
		store.Inc(1)
		store.Inc(2)
		store.Inc(2)
		store.Inc(3)
		key := store.LFU()
		require.Equal(t, 3, key)
	})
}
