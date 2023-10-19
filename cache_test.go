package memcache_test

import (
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

	t.Run("returns an error when an option returns an error", func(t *testing.T) {
		t.Parallel()

		var errDummy error = errors.New("dummy")
		c, err := memcache.New[int, string](func(c *memcache.Cache[int, string]) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})

	t.Run("enables passive expiration", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.New[int, string](memcache.WithPassiveExpiration[int, string]())
		require.NoError(t, err)
		require.True(t, c.GetExpireOnGet())
	})
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns value of key & true when key exists in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		unlock()

		got, ok := c.Get(1)
		require.True(t, ok)
		require.Equal(t, "a", got)
	})

	t.Run("returns empty string & false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()

		got, ok := c.Get(1)
		require.False(t, ok)
		require.Equal(t, "", got)
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.New[int, string](memcache.WithPassiveExpiration[int, string]())
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a", ExpireAt: &expireAt}
		unlock()

		got, ok := c.Get(1)
		require.False(t, ok)
		require.Equal(t, "", got)
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a", ExpireAt: &expireAt}
		unlock()

		got, ok := c.Get(1)
		require.True(t, ok)
		require.Equal(t, "a", got)
	})
}

func TestCache_Has(t *testing.T) {
	t.Parallel()

	t.Run("returns true if key exists in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		unlock()

		require.True(t, c.Has(1))
	})

	t.Run("returns false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()

		require.False(t, c.Has(1))
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.New[int, string](memcache.WithPassiveExpiration[int, string]())
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a", ExpireAt: &expireAt}
		unlock()

		require.False(t, c.Has(1))
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a", ExpireAt: &expireAt}
		unlock()

		require.True(t, c.Has(1))
	})
}

func TestCache_Set(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache at provided key", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()

		c.Set(1, "a")

		store, unlock := c.GetStore()
		defer unlock()
		require.Contains(t, store, 1)
		require.Equal(t, "a", store[1].Value)
	})
}

func TestCache_SetEx(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache with a TTL", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()

		c.SetEx(1, "a", 1*time.Minute)

		store, unlock := c.GetStore()
		defer unlock()
		require.Contains(t, store, 1)
		require.Equal(t, "a", store[1].Value)
		require.Greater(t, *store[1].ExpireAt, time.Now())
	})
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes value from the cache at provided key", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		unlock()

		c.Delete(1)

		store, unlock = c.GetStore()
		defer unlock()
		require.NotContains(t, store, 1)
	})
}

func TestCache_Flush(t *testing.T) {
	t.Parallel()

	t.Run("successfully flushes the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		store[2] = memcache.Item[int, string]{Value: "b"}
		store[3] = memcache.Item[int, string]{Value: "c"}
		unlock()

		c.Flush()

		store, unlock = c.GetStore()
		defer unlock()
		require.Empty(t, store)
	})
}

func TestCache_Size(t *testing.T) {
	t.Parallel()

	t.Run("returns the size of the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		store[2] = memcache.Item[int, string]{Value: "b"}
		store[3] = memcache.Item[int, string]{Value: "c"}
		unlock()

		require.Equal(t, 3, c.Size())
	})
}

func TestCache_Keys(t *testing.T) {
	t.Parallel()

	t.Run("returns the keys of the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.New[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		store[2] = memcache.Item[int, string]{Value: "b"}
		store[3] = memcache.Item[int, string]{Value: "c"}
		unlock()

		require.ElementsMatch(t, []int{1, 2, 3}, c.Keys())
	})
}
