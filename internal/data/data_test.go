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

func TestItem_TTL(t *testing.T) {
	t.Parallel()

	t.Run("returns time remaining when the item is not expired", func(t *testing.T) {
		t.Parallel()

		ttl := 1 * time.Minute
		now := time.Now().Add(1 * time.Minute)
		i := data.Item[int, string]{ExpireAt: &now}
		require.Greater(t, *i.TTL(), ttl-1*time.Second)
	})

	t.Run("returns nil when the item has no expiry", func(t *testing.T) {
		t.Parallel()

		i := data.Item[int, string]{ExpireAt: nil}
		require.Nil(t, i.TTL())
	})

	t.Run("returns 0 when the item is expired", func(t *testing.T) {
		t.Parallel()

		now := time.Now().Add(-1 * time.Minute)
		i := data.Item[int, string]{ExpireAt: &now}
		require.Equal(t, time.Duration(0), *i.TTL())
	})
}
