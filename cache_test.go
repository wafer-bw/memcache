package memcache_test

import (
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a new cache", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.Open[int, string]()
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			c, err := memcache.Open[int, string](nil, nil)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	})

	t.Run("returns an error when an option returns an error", func(t *testing.T) {
		t.Parallel()

		var errDummy error = errors.New("dummy")
		c, err := memcache.Open[int, string](func(c *memcache.Cache[int, string]) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})

	t.Run("with passive expiration enables passive expiration", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.Open[int, string](memcache.WithPassiveExpiration[int, string]())
		require.NoError(t, err)
		require.True(t, c.GetExpireOnGet())
	})

	t.Run("with expirer sets and runs the expirer", func(t *testing.T) {
		t.Parallel()

		ran := new(bool)
		interval := 25 * time.Millisecond
		expirer := func(store map[int]memcache.Item[int, string]) {
			*ran = true
		}

		c, _ := memcache.Open[int, string](memcache.WithActiveExpiration[int, string](expirer, interval))
		defer c.Close()
		time.Sleep(interval * 2)
		require.NotNil(t, c.GetExpirer())

		// increase chances of race condition to make test more reliable
		for i := 0; i < 100000; i++ {
			c.RLock()
			require.True(t, *ran)
			c.RUnlock()
		}
	})

	t.Run("with expirer returns an error if the expirer function is nil", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.Open[int, int](memcache.WithActiveExpiration[int, int](nil, 1*time.Second))
		require.ErrorIs(t, err, memcache.ErrNilExpirerFunc)
	})

	t.Run("with expirer returns an error if the interval is less than or equal to 0", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.Open[int, int](memcache.WithActiveExpiration[int, int](memcache.DeleteAllExpiredKeys, 0))
		require.ErrorIs(t, err, memcache.ErrInvalidInterval)
	})
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns value of key & true when key exists in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		unlock()

		got, ok := c.Get(1)
		require.True(t, ok)
		require.Equal(t, "a", got)
	})

	t.Run("returns empty string & false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()

		got, ok := c.Get(1)
		require.False(t, ok)
		require.Equal(t, "", got)
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.Open[int, string](memcache.WithPassiveExpiration[int, string]())
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
		c, _ := memcache.Open[int, string]()
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

		c, _ := memcache.Open[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		unlock()

		require.True(t, c.Has(1))
	})

	t.Run("returns false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()

		require.False(t, c.Has(1))
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.Open[int, string](memcache.WithPassiveExpiration[int, string]())
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a", ExpireAt: &expireAt}
		unlock()

		require.False(t, c.Has(1))
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()

		expireAt := time.Now()
		c, _ := memcache.Open[int, string]()
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

		c, _ := memcache.Open[int, string]()

		c.Set(1, "a")

		store, unlock := c.GetStore()
		defer unlock()
		require.Contains(t, store, 1)
		require.Equal(t, "a", store[1].Value)
	})

	t.Run("demonstrates unsafe usage of pointer values stored in cache", func(t *testing.T) {
		t.Parallel()

		v := false
		c, _ := memcache.Open[int, *bool]()

		c.Set(1, &v)
		v = true

		store, unlock := c.GetStore()
		defer unlock()
		require.Contains(t, store, 1)
		require.Equal(t, true, *store[1].Value)
	})
}

func TestCache_SetEx(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache with a TTL", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()

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

	t.Run("successfully deletes key from cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		unlock()

		c.Delete(1)

		store, unlock = c.GetStore()
		defer unlock()
		require.NotContains(t, store, 1)
	})

	t.Run("successfully deletes keys from cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		store[2] = memcache.Item[int, string]{Value: "b"}
		unlock()

		c.Delete(1, 2)

		store, unlock = c.GetStore()
		defer unlock()
		require.NotContains(t, store, 1)
		require.NotContains(t, store, 2)
	})
}

func TestCache_Flush(t *testing.T) {
	t.Parallel()

	t.Run("successfully flushes the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()
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

		c, _ := memcache.Open[int, string]()
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

		c, _ := memcache.Open[int, string]()
		store, unlock := c.GetStore()
		store[1] = memcache.Item[int, string]{Value: "a"}
		store[2] = memcache.Item[int, string]{Value: "b"}
		store[3] = memcache.Item[int, string]{Value: "c"}
		unlock()

		require.ElementsMatch(t, []int{1, 2, 3}, c.Keys())
	})
}

func TestCache_Close(t *testing.T) {
	t.Parallel()

	t.Run("successfully closes the cache", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()

		c.Close()
		require.True(t, c.Closed())
	})

	t.Run("subsequent calls to close do not panic", func(t *testing.T) {
		t.Parallel()

		c, _ := memcache.Open[int, string]()

		c.Close()
		require.NotPanics(t, func() { c.Close() })
	})

	t.Run("cache with no goroutines is garbage collected after releasing without closing", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})

		cache := func() *memcache.Cache[int, string] {
			c, _ := memcache.Open[int, string]()
			runtime.SetFinalizer(c, func(_ *memcache.Cache[int, string]) {
				close(ch)
			})
			return c
		}()

		cache.Flush() // use the cache once
		cache = nil   // release the cache
		runtime.GC()  // explicitly run garbage collection

		select {
		case <-time.After(250 * time.Millisecond):
			t.Fatal("cache was not garbage collected")
		case <-ch:
			// cache was garbage collected
		}
	})

	t.Run("cache with goroutines is garbage collected after releasing & closing", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})

		cache := func() *memcache.Cache[int, string] {
			interval := 1 * time.Second
			expirer := memcache.DeleteAllExpiredKeys[int, string]
			c, _ := memcache.Open[int, string](memcache.WithActiveExpiration(expirer, interval))
			runtime.SetFinalizer(c, func(_ *memcache.Cache[int, string]) {
				close(ch)
			})
			return c
		}()

		cache.Flush() // use the cache once
		cache.Close() // close the cache
		cache = nil   // release the cache
		runtime.GC()  // explicitly run garbage collection

		select {
		case <-time.After(250 * time.Millisecond):
			t.Fatal("cache was not garbage collected")
		case <-ch:
			// cache was garbage collected
		}
	})
}
