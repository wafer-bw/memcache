package memcache_test

import (
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
)

const cacheSize = 100

var policies = map[string]func(size int, options ...memcache.Option[int, int]) (*memcache.Cache[int, int], error){
	"noevict": func(_ int, options ...memcache.Option[int, int]) (*memcache.Cache[int, int], error) {
		return memcache.Open[int, int](options...)
	},
	"lru": func(size int, options ...memcache.Option[int, int]) (*memcache.Cache[int, int], error) {
		options = append(options, memcache.WithLRUEviction[int, int](size))
		return memcache.Open[int, int](options...)
	},
}

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
		require.True(t, c.PassiveExpiration())
	})

	t.Run("with expirer sets and runs the expirer", func(t *testing.T) {
		t.Parallel()

		// TODO: once expirer func is an interface use a mock to expect call.

		ch := make(chan struct{})
		interval := 25 * time.Millisecond
		expirer := func(items map[int]memcache.Item[int, string]) {
			close(ch)
		}

		c, _ := memcache.Open[int, string](memcache.WithActiveExpiration[int, string](expirer, interval))
		defer c.Close()

		select {
		case <-time.After(interval * 2):
			t.Fatal("expirer was not called")
		case <-ch:
			// expirer was called
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

	t.Run("with lru eviction sets the store to an lru store", func(t *testing.T) {
		t.Parallel()
		c, _ := memcache.Open[int, string](memcache.WithLRUEviction[int, string](2))
		store := c.Store()

		expected := memcache.LRUStore[int, string]{}.Underlying
		require.IsType(t, expected, store)
	})

	t.Run("with lru eviction returns an error if the capacity is less than or equal to 1", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.Open[int, int](memcache.WithLRUEviction[int, int](1))
		require.ErrorIs(t, err, memcache.ErrInvalidCapacity)
	})
}

func TestCache_Set(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache at provided key", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)

				c.Set(1, 1)
				items, unlock := c.Store().Items()
				defer unlock()
				require.Contains(t, items, 1)
				require.Equal(t, 1, items[1].Value)
			})
		}
	})

	t.Run("demonstrates unsafe usage of pointer values stored in cache", func(t *testing.T) {
		t.Parallel()

		v := false
		c, _ := memcache.Open[int, *bool]()

		c.Set(1, &v)
		v = true

		items, unlock := c.Store().Items()
		defer unlock()
		require.Contains(t, items, 1)
		require.Equal(t, true, *items[1].Value)
	})
}

func TestCache_SetEx(t *testing.T) {
	t.Parallel()

	t.Run("successfully stores value in the cache with a TTL", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)

				c.SetEx(1, 1, 1*time.Minute)
				items, unlock := c.Store().Items()
				defer unlock()
				require.Contains(t, items, 1)
				require.Equal(t, 1, items[1].Value)
				require.Greater(t, *items[1].ExpireAt, time.Now())
			})
		}
	})
}

func TestCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns value of key & true when key exists in the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})

				got, ok := c.Get(1)
				require.True(t, ok)
				require.Equal(t, 1, got)
			})
		}
	})

	t.Run("returns empty string & false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)

				got, ok := c.Get(1)
				require.False(t, ok)
				require.Equal(t, 0, got)
			})
		}
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				expireAt := time.Now()
				c, _ := newCache(cacheSize, memcache.WithPassiveExpiration[int, int]())
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				got, ok := c.Get(1)
				require.False(t, ok)
				require.Equal(t, 0, got)
			})
		}
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				expireAt := time.Now()
				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				got, ok := c.Get(1)
				require.True(t, ok)
				require.Equal(t, 1, got)
			})
		}
	})
}

func TestCache_Has(t *testing.T) {
	t.Parallel()

	t.Run("returns true if key exists in the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})

				ok := c.Has(1)
				require.True(t, ok)
			})
		}
	})

	t.Run("returns false when key does not exist in the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)

				ok := c.Has(1)
				require.False(t, ok)
			})
		}
	})

	t.Run("deletes expired keys when passive expiration is enabled", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				expireAt := time.Now()
				c, _ := newCache(cacheSize, memcache.WithPassiveExpiration[int, int]())
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				ok := c.Has(1)
				require.False(t, ok)
			})
		}
	})

	t.Run("does not delete expired keys when passive expiration is disabled", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				expireAt := time.Now()
				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				ok := c.Has(1)
				require.True(t, ok)
			})
		}
	})
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes key from cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})

				c.Delete(1)
				items, unlock := store.Items()
				defer unlock()
				require.NotContains(t, items, 1)
			})
		}
	})

	t.Run("successfully deletes keys from cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})
				store.Set(2, memcache.Item[int, int]{Value: 2})

				c.Delete(1, 2)
				items, unlock := store.Items()
				defer unlock()
				require.NotContains(t, items, 1)
				require.NotContains(t, items, 2)
			})
		}
	})

	t.Run("does not panic when deleting keys that do not exist in the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)

				require.NotPanics(t, func() {
					c.Delete(1, 2)
				})
			})
		}
	})
}

func TestCache_Flush(t *testing.T) {
	t.Parallel()

	t.Run("successfully flushes the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})
				store.Set(2, memcache.Item[int, int]{Value: 2})
				store.Set(3, memcache.Item[int, int]{Value: 3})

				c.Flush()
				items, unlock := store.Items()
				defer unlock()
				require.Empty(t, items)
			})
		}
	})
}

func TestCache_Size(t *testing.T) {
	t.Parallel()

	t.Run("returns the size of the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})
				store.Set(2, memcache.Item[int, int]{Value: 2})
				store.Set(3, memcache.Item[int, int]{Value: 3})

				size := c.Size()
				require.Equal(t, 3, size)
			})
		}
	})
}

func TestCache_Keys(t *testing.T) {
	t.Parallel()

	t.Run("returns the keys of the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, memcache.Item[int, int]{Value: 1})
				store.Set(2, memcache.Item[int, int]{Value: 1})
				store.Set(3, memcache.Item[int, int]{Value: 1})

				keys := c.Keys()
				require.ElementsMatch(t, []int{1, 2, 3}, keys)
			})
		}
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
