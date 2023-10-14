package memcache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
	"github.com/wafer-bw/memcache/internal/record"
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

	t.Run("executes provided options", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var a, b *int = new(int), new(int)
		c, err := memcache.New[int, string](ctx,
			func(_ *memcache.CacheConfig) error {
				*a = 1
				return nil
			},
			func(_ *memcache.CacheConfig) error {
				*b = 2
				return nil
			},
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, a)
		require.Equal(t, 1, *a)
		require.NotNil(t, b)
		require.Equal(t, 2, *b)
	})

	t.Run("returns an error when an option returns an error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var errDummy error = errors.New("dummy")
		c, err := memcache.New[int, string](ctx, func(c *memcache.CacheConfig) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})

	t.Run("sets expire on get to true", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, err := memcache.New[int, string](ctx, memcache.WithExpireOnGet())
		require.NoError(t, err)
		require.True(t, c.GetExpireOnGet())
	})

	t.Run("sets expiration interval", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		i := 1 * time.Second
		c, err := memcache.New[int, string](ctx, memcache.WithExpirationInterval(i))
		require.NoError(t, err)
		require.Equal(t, i, c.GetExpirationInterval())
	})
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns value that exists in the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		k, v := 1, "a"
		c, _ := memcache.New[int, string](ctx)
		store := c.GetStore()
		store[k] = record.Record[string]{Value: v}

		got, ok := c.Get(k)
		require.True(t, ok)
		require.Equal(t, v, got)
	})

	t.Run("returns false when value does not exist in the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)

		_, ok := c.Get(1)
		require.False(t, ok)
	})

	t.Run("expires stale values when expire on get is on", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		now := time.Now()
		k, v := 1, "a"
		c, _ := memcache.New[int, string](ctx, memcache.WithExpireOnGet())
		store := c.GetStore()
		store[k] = record.Record[string]{Value: v, ExpireAt: &now}

		_, ok := c.Get(k)
		require.False(t, ok)
	})

	t.Run("does not expire stale values when expire on get is off", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		now := time.Now()
		k, v := 1, "a"
		c, _ := memcache.New[int, string](ctx)
		store := c.GetStore()
		store[k] = record.Record[string]{Value: v, ExpireAt: &now}

		_, ok := c.Get(k)
		require.True(t, ok)
	})
}

func TestCache_Set(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache at provided key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		k, v := 1, "a"
		c, _ := memcache.New[int, string](ctx)
		store := c.GetStore()

		c.Set(k, v)
		require.Contains(t, store, k)
		require.Equal(t, v, store[k].Value)
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		require.NotPanics(t, func() {
			c.Set(1, "a", nil, nil)
		})
	})

	t.Run("executes provided options", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var a, b *int = new(int), new(int)
		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a",
			func(_ *memcache.ValueConfig) {
				*a = 1
			},
			func(_ *memcache.ValueConfig) {
				*b = 2
			},
		)
		require.NotNil(t, a)
		require.Equal(t, 1, *a)
		require.NotNil(t, b)
		require.Equal(t, 2, *b)
	})

	t.Run("successfully stores value in the cache with a TTL", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		c.Set(1, "a", memcache.WithTTL(1*time.Minute))
		require.NotNil(t, c.GetStore()[1].ExpireAt)
	})
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes value from the cache at provided key", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		k, v := 1, "a"
		c, _ := memcache.New[int, string](ctx)
		store := c.GetStore()
		store[k] = record.Record[string]{Value: v}

		c.Delete(k)
		require.NotContains(t, store, k)
	})
}

func TestCache_Flush(t *testing.T) {
	t.Parallel()

	t.Run("successfully flushes the cache", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		c, _ := memcache.New[int, string](ctx)
		store := c.GetStore()
		store[1] = record.Record[string]{Value: "a"}
		store[2] = record.Record[string]{Value: "b"}
		store[3] = record.Record[string]{Value: "c"}

		c.Flush()
		require.Empty(t, store)
	})
}
