package data_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/data"
)

func TestItem_IsExpired(t *testing.T) {
	t.Parallel()

	t.Run("returns true when ExpireAt is in the past", func(t *testing.T) {
		t.Parallel()

		now := time.Now().Add(-1 * time.Minute)
		i := data.Item[int, string]{ExpireAt: &now}
		require.True(t, i.IsExpired())
	})

	t.Run("returns false when ExpireAt is nil", func(t *testing.T) {
		t.Parallel()

		var i data.Item[int, string]
		require.False(t, i.IsExpired())
	})

	t.Run("returns false when ExpireAt is in the future", func(t *testing.T) {
		t.Parallel()

		now := time.Now().Add(1 * time.Minute)
		i := data.Item[int, string]{ExpireAt: &now}
		require.False(t, i.IsExpired())
	})
}
