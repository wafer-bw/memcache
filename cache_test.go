package memcache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a new cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, err := memcache.New[int, string](ctx)
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		require.NotPanics(t, func() {
			c, err := memcache.New[int, string](ctx, nil, nil)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	})

	t.Run("returns an error when an option returns an error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var errDummy error = errors.New("dummy")
		c, err := memcache.New[int, string](ctx, func(c *memcache.CacheConfig) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})

	t.Run("enables passive expiration", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, err := memcache.New[int, string](ctx, memcache.WithPassiveExpiration())
		require.NoError(t, err)
		require.True(t, c.GetExpireOnGet())
	})
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns value of key & true when key exists in the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a")
		got, ok := c.Get(1)
		require.True(t, ok)
		require.Equal(t, "a", got)
	})

	t.Run("returns empty string & false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		got, ok := c.Get(1)
		require.False(t, ok)
		require.Equal(t, "", got)
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx, memcache.WithPassiveExpiration())
		c.Set(1, "a", memcache.WithTTL(0*time.Second))
		got, ok := c.Get(1)
		require.False(t, ok)
		require.Equal(t, "", got)
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a", memcache.WithTTL(0*time.Second))
		got, ok := c.Get(1)
		require.True(t, ok)
		require.Equal(t, "a", got)
	})
}

func TestCache_Has(t *testing.T) {
	t.Parallel()

	t.Run("returns true if key exists in the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a")
		require.True(t, c.Has(1))
	})

	t.Run("returns false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		require.False(t, c.Has(1))
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx, memcache.WithPassiveExpiration())
		c.Set(1, "a", memcache.WithTTL(0*time.Second))
		require.False(t, c.Has(1))
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a", memcache.WithTTL(0*time.Second))
		require.True(t, c.Has(1))
	})
}

func TestCache_Set(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache at provided key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		k, v := 1, "a"
		c, _ := memcache.New[int, string](ctx)

		c.Set(k, v)

		got, ok := c.Get(k)
		require.True(t, ok)
		require.Equal(t, v, got)
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		require.NotPanics(t, func() {
			c.Set(1, "a", nil, nil)
		})
	})

	t.Run("successfully stores value in the cache with a TTL", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a", memcache.WithTTL(1*time.Minute))
		store, unlock := c.GetStore()
		defer unlock()
		require.NotNil(t, store[1].ExpireAt)
	})
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes value from the cache at provided key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a")
		c.Delete(1)
		require.False(t, c.Has(1))
	})
}

func TestCache_Flush(t *testing.T) {
	t.Parallel()

	t.Run("successfully flushes the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a")
		c.Set(2, "b")
		c.Set(3, "c")

		c.Flush()

		require.Equal(t, 0, c.Length())
	})
}

func TestCache_Length(t *testing.T) {
	t.Parallel()

	t.Run("returns the length of the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a")
		c.Set(2, "b")
		c.Set(3, "c")
		require.Equal(t, 3, c.Length())
	})
}

func TestCache_Keys(t *testing.T) {
	t.Parallel()

	t.Run("returns the keys of the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a")
		c.Set(2, "b")
		c.Set(3, "c")
		require.ElementsMatch(t, []int{1, 2, 3}, c.Keys())
	})
}
