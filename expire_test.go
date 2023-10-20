package memcache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

func TestDeleteAllExpiredKeys(t *testing.T) {
	t.Parallel()

	t.Run("deletes all expired keys", func(t *testing.T) {
		now := time.Now()
		store := map[int]memcache.Item[int, string]{
			1: {Value: "a"},
			2: {Value: "b", ExpireAt: &now},
			3: {Value: "c", ExpireAt: &now},
		}

		memcache.DeleteAllExpiredKeys[int, string](store)
		require.Len(t, store, 1)
	})
}
