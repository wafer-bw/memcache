package memcache_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
	"github.com/wafer-bw/memcache/internal/record"
)

func TestFullScanExpirer(t *testing.T) {
	t.Parallel()

	t.Run("deletes all expired keys from the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		now := time.Now()

		c, _ := memcache.New[int, int](ctx, memcache.WithExpirationInterval(1*time.Millisecond))
		store := c.GetStore()

		store[1] = record.Record[int]{Value: 1, ExpireAt: &now}
		store[2] = record.Record[int]{Value: 2, ExpireAt: &now}
		store[3] = record.Record[int]{Value: 3, ExpireAt: &now}

		time.Sleep(2 * time.Millisecond)
		require.NotContains(t, store, 1)
		require.NotContains(t, store, 2)
		require.NotContains(t, store, 3)
	})
}
