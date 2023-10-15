package memcache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

func TestItem_IsExpired(t *testing.T) {
	t.Parallel()

	t.Run("returns false when ExpireAt is nil", func(t *testing.T) {
		t.Parallel()

		var r memcache.Item[string]
		require.False(t, r.IsExpired())
	})

	t.Run("returns false when ExpireAt is in the future", func(t *testing.T) {
		t.Parallel()

		now := time.Now().Add(1 * time.Minute)
		r := memcache.Item[string]{ExpireAt: &now}
		require.False(t, r.IsExpired())
	})

	t.Run("returns true when ExpireAt is in the past", func(t *testing.T) {
		t.Parallel()

		now := time.Now().Add(-1 * time.Minute)
		r := memcache.Item[string]{ExpireAt: &now}
		require.True(t, r.IsExpired())
	})
}
