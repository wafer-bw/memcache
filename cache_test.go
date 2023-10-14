package memcache_test

import (
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

		c, err := memcache.New[int, string]()
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			c, err := memcache.New[int, string](nil, nil)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	})

	t.Run("executes provided options", func(t *testing.T) {
		t.Parallel()

		var a, b *int = new(int), new(int)
		c, err := memcache.New[int, string](
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

		var errDummy error = errors.New("dummy")
		c, err := memcache.New[int, string](func(c *memcache.CacheConfig) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns value that exists in the cache", func(t *testing.T) {
		t.Parallel()

		k, v := 1, "a"
		c, _ := memcache.New[int, string]()
		store := c.GetStore()
		store[k] = record.Record[string]{Value: v}

		got, ok := c.Get(k)
		require.True(t, ok)
		require.Equal(t, v, got)
	})

	t.Run("returns false when value does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()

		_, ok := c.Get(1)
		require.False(t, ok)
	})
}

func TestCache_Set(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache at provided key", func(t *testing.T) {
		t.Parallel()

		k, v := 1, "a"
		c, _ := memcache.New[int, string]()
		store := c.GetStore()

		c.Set(k, v)
		require.Contains(t, store, k)
		require.Equal(t, v, store[k].Value)
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		require.NotPanics(t, func() {
			c.Set(1, "a", nil, nil)
		})
	})

	t.Run("executes provided options", func(t *testing.T) {
		t.Parallel()

		var a, b *int = new(int), new(int)
		c, _ := memcache.New[int, string]()
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
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes value from the cache at provided key", func(t *testing.T) {
		t.Parallel()

		k, v := 1, "a"
		c, _ := memcache.New[int, string]()
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

		c, _ := memcache.New[int, string]()
		store := c.GetStore()
		store[1] = record.Record[string]{Value: "a"}
		store[2] = record.Record[string]{Value: "b"}
		store[3] = record.Record[string]{Value: "c"}

		c.Flush()
		require.Empty(t, store)
	})
}

func TestWithEvictionInterval(t *testing.T) {
	t.Parallel()

	t.Run("sets the eviction interval", func(t *testing.T) {
		t.Parallel()

		i := 1 * time.Second
		c := memcache.CacheConfig{}
		err := memcache.WithEvictionInterval(i)(&c)
		require.NoError(t, err)
		require.Equal(t, i, c.GetEvictionInterval())
	})
}

func TestWithExpirationInterval(t *testing.T) {
	t.Parallel()

	t.Run("sets the expiration interval", func(t *testing.T) {
		t.Parallel()

		i := 1 * time.Second
		c := memcache.CacheConfig{}
		err := memcache.WithExpirationInterval(i)(&c)
		require.NoError(t, err)
		require.Equal(t, i, c.GetExpirationInterval())
	})
}

func TestWithTTL(t *testing.T) {
	t.Parallel()

	t.Run("sets the TTL", func(t *testing.T) {
		t.Parallel()

		d := 1 * time.Minute
		c := memcache.ValueConfig{}
		memcache.WithTTL(d)(&c)
		got := c.GetExpireAt()
		expect := time.Now().Add(d).Truncate(time.Second).UnixNano()
		require.NotNil(t, got)
		require.Equal(t, expect, got.Truncate(time.Second).UnixNano())
	})
}
