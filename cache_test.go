package memcache_test

import (
	"errors"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/store/lru"
	"github.com/wafer-bw/memcache/internal/store/noevict"
)

type cacher[K comparable, V any] interface {
	Set(key K, value V)
	SetEx(key K, value V, ttl time.Duration)
	Get(key K) (V, bool)
	TTL(key K) (*time.Duration, bool)
	Delete(keys ...K)
	Size() int
	Keys() []K
	Flush()
	Close()

	// TODO - add the following methods:
	// Need:
	// - Scan()    // iterate over keys in cache (requires upcoming go iterators).
	// Maybe:
	// - Random()  // return random key/value from cache.
	// - Persist() // remove ttl from key.
	// - Expire()  // set ttl for key.
}

var _ cacher[int, int] = (*memcache.Cache[int, int])(nil)
var _ memcache.Storer[int, int] = (*lru.Store[int, int])(nil)
var _ memcache.Storer[int, int] = (*noevict.Store[int, int])(nil)

const cacheSize = 100

var policies = map[string]func(size int, options ...memcache.Option[int, int]) (*memcache.Cache[int, int], error){
	noevict.PolicyName: func(_ int, options ...memcache.Option[int, int]) (*memcache.Cache[int, int], error) {
		return memcache.OpenNoEvictionCache[int, int](options...)
	},
	lru.PolicyName: func(size int, options ...memcache.Option[int, int]) (*memcache.Cache[int, int], error) {
		return memcache.OpenLRUCache[int, int](size, options...)
	},
}

func TestInvalidCapacityError_Error(t *testing.T) {
	t.Parallel()

	t.Run("returns error message", func(t *testing.T) {
		t.Parallel()

		err := memcache.InvalidCapacityError{Capacity: 0, Minimum: 1, Policy: "active"}
		require.Equal(t, "capacity 0 must be greater than 1 for active caches", err.Error())
	})
}

func TestCache_concurrentAccess(t *testing.T) {
	t.Run("passive expiration disabled", func(t *testing.T) {
		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				n := 1000
				cache, _ := newCache(n)

				for i := 0; i < n; i++ {
					cache.Set(i, i)
				}

				var wg sync.WaitGroup
				for i := 0; i < n; i++ {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
						_, _ = cache.Get(i)
						cache.Set(i, i*2)
					}(i)
				}
				wg.Wait()

				for i := 0; i < n; i++ {
					value, ok := cache.Get(i)
					require.True(t, ok, i)
					require.Equal(t, i*2, value)
				}
			})
		}
	})

	t.Run("passive expiration enabled", func(t *testing.T) {
		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				n := 1000
				n2 := n * 2
				cache, _ := newCache(n, memcache.WithPassiveExpiration[int, int]())

				for i := 0; i < n; i++ {
					cache.SetEx(i, i, 1*time.Millisecond)
				}

				var wg sync.WaitGroup
				for i := 0; i < n2; i++ {
					wg.Add(1)
					go func() {
						v := rand.Intn(n2 - 1)
						defer wg.Done()
						cache.Get(v)
						cache.SetEx(v, 1, 1*time.Millisecond)
					}()
				}
				wg.Wait()
			})
		}
	})
}

func TestCache_activeExpiration(t *testing.T) {
	t.Parallel()

	t.Run("deletes expired keys", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				ttl := 1 * time.Millisecond

				cache, _ := newCache(cacheSize, memcache.WithActiveExpiration[int, int](ttl))
				defer cache.Close()

				cache.SetEx(1, 1, ttl)
				cache.SetEx(2, 2, ttl)
				cache.SetEx(3, 3, ttl)

				time.Sleep(5 * ttl)

				items, unlock := cache.Store().Items()
				defer unlock()
				require.Empty(t, items)
			})
		}
	})
}

func TestOpenNoEvictionCache(t *testing.T) {
	t.Parallel()

	t.Run("returns a new no eviction cache", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.OpenNoEvictionCache[int, string]()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.IsType(t, &noevict.Store[int, string]{}, c.Store())
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			c, err := memcache.OpenNoEvictionCache[int, string](nil, nil)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	})

	t.Run("returns an error when an option returns an error", func(t *testing.T) {
		t.Parallel()

		errDummy := errors.New("dummy")

		c, err := memcache.OpenNoEvictionCache[int, string](func(c *memcache.Cache[int, string]) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})

	t.Run("returns an error when opening the store returns an error", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.OpenNoEvictionCache[int, string](memcache.WithCapacity[int, string](-1))
		require.Error(t, err)
		require.Nil(t, c)
	})

	t.Run("with passive expiration enables passive expiration", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.OpenNoEvictionCache[int, string](memcache.WithPassiveExpiration[int, string]())
		require.NoError(t, err)
		require.True(t, c.PassiveExpiration())
	})

	t.Run("with active expiration enables active expiration", func(t *testing.T) {
		t.Parallel()

		interval := 25 * time.Millisecond

		c, err := memcache.OpenNoEvictionCache[int, string](memcache.WithActiveExpiration[int, string](interval))
		require.NoError(t, err)
		defer c.Close()
		require.Equal(t, interval, c.ExpirationInterval())
	})

	t.Run("with active expiration returns an error if the interval is less than or equal to 0", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.OpenNoEvictionCache[int, int](memcache.WithActiveExpiration[int, int](0 * time.Second))
		require.ErrorIs(t, err, memcache.ErrInvalidInterval)
	})

	t.Run("with capacity sets capacity", func(t *testing.T) {
		t.Parallel()

		capacity := 10

		c, err := memcache.OpenNoEvictionCache[int, string](memcache.WithCapacity[int, string](capacity))
		require.NoError(t, err)
		require.Equal(t, capacity, c.Capacity())
	})

	t.Run("returns an error if capacity is less than 0", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.OpenNoEvictionCache[int, int](memcache.WithCapacity[int, int](-1))
		require.ErrorAs(t, err, &memcache.InvalidCapacityError{})
	})
}

func TestOpenLRUCache(t *testing.T) {
	t.Parallel()

	t.Run("returns a new lru cache", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.OpenLRUCache[int, string](10)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.IsType(t, &lru.Store[int, string]{}, c.Store())
	})

	t.Run("does not panic when provided nil options", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			c, err := memcache.OpenLRUCache[int, string](10, nil, nil)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	})

	t.Run("returns an error when an option returns an error", func(t *testing.T) {
		t.Parallel()

		errDummy := errors.New("dummy")

		c, err := memcache.OpenLRUCache[int, string](10, func(c *memcache.Cache[int, string]) error { return errDummy })
		require.ErrorIs(t, err, errDummy)
		require.Nil(t, c)
	})

	t.Run("returns an error when opening the store returns an error", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.OpenLRUCache[int, string](0)
		require.Error(t, err)
		require.Nil(t, c)
	})

	t.Run("with passive expiration enables passive expiration", func(t *testing.T) {
		t.Parallel()

		c, err := memcache.OpenLRUCache[int, string](10, memcache.WithPassiveExpiration[int, string]())
		require.NoError(t, err)
		require.True(t, c.PassiveExpiration())
	})

	t.Run("with active expiration enables active expiration", func(t *testing.T) {
		t.Parallel()

		interval := 25 * time.Millisecond

		c, err := memcache.OpenLRUCache[int, string](10, memcache.WithActiveExpiration[int, string](interval))
		require.NoError(t, err)
		defer c.Close()
		require.Equal(t, interval, c.ExpirationInterval())
	})

	t.Run("with active expiration returns an error if the interval is less than or equal to 0", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.OpenLRUCache[int, int](10, memcache.WithActiveExpiration[int, int](0*time.Second))
		require.ErrorIs(t, err, memcache.ErrInvalidInterval)
	})

	t.Run("returns an error if the capacity is less than 1", func(t *testing.T) {
		t.Parallel()

		_, err := memcache.OpenLRUCache[int, int](0)
		require.ErrorAs(t, err, &memcache.InvalidCapacityError{})
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
				store.Set(1, data.Item[int, int]{Value: 1})

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

	t.Run("returns empty string & false when expired key exists in the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				expireAt := time.Now()
				c, _ := newCache(cacheSize)
				store := c.Store()
				store.Set(1, data.Item[int, int]{Value: 1, ExpireAt: &expireAt})

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
				store.Set(1, data.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				_, ok := c.Get(1)
				require.False(t, ok)

				items, unlock := store.Items()
				defer unlock()
				require.NotContains(t, items, 1)
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
				store.Set(1, data.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				_, ok := c.Get(1)
				require.False(t, ok)

				items, unlock := store.Items()
				defer unlock()
				require.Contains(t, items, 1)
			})
		}
	})
}

func TestCache_TTL(t *testing.T) {
	t.Parallel()

	t.Run("returns remaining ttl of expiring item", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()

				setTTL := 1 * time.Minute
				expireAt := time.Now().Add(setTTL)
				store.Set(1, data.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				ttl, ok := c.TTL(1)
				require.True(t, ok)
				require.Greater(t, *ttl, setTTL-1*time.Second)
			})
		}
	})

	t.Run("returns zero ttl for expired item", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()

				setTTL := -1 * time.Minute
				expireAt := time.Now().Add(setTTL)
				store.Set(1, data.Item[int, int]{Value: 1, ExpireAt: &expireAt})

				ttl, ok := c.TTL(1)
				require.True(t, ok)
				require.Equal(t, time.Duration(0), *ttl)
			})
		}
	})

	t.Run("returns nil ttl for non-expiring item", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				store := c.Store()

				store.Set(1, data.Item[int, int]{Value: 1})

				ttl, ok := c.TTL(1)
				require.True(t, ok)
				require.Nil(t, ttl)
			})
		}
	})

	t.Run("returns nil ttl and false if key does not exist", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)

				ttl, ok := c.TTL(1)
				require.False(t, ok)
				require.Nil(t, ttl)
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
				store.Set(1, data.Item[int, int]{Value: 1})

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
				store.Set(1, data.Item[int, int]{Value: 1})
				store.Set(2, data.Item[int, int]{Value: 2})

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
				store.Set(1, data.Item[int, int]{Value: 1})
				store.Set(2, data.Item[int, int]{Value: 2})
				store.Set(3, data.Item[int, int]{Value: 3})

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
				store.Set(1, data.Item[int, int]{Value: 1})
				store.Set(2, data.Item[int, int]{Value: 2})
				store.Set(3, data.Item[int, int]{Value: 3})

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
				store.Set(1, data.Item[int, int]{Value: 1})
				store.Set(2, data.Item[int, int]{Value: 1})
				store.Set(3, data.Item[int, int]{Value: 1})

				keys := c.Keys()
				require.Contains(t, keys, 1)
				require.Contains(t, keys, 2)
				require.Contains(t, keys, 3)
			})
		}
	})
}

func TestCache_Close(t *testing.T) {
	t.Parallel()

	t.Run("successfully closes the cache", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				c.Close()
				require.True(t, c.Closed())
			})
		}
	})

	t.Run("subsequent calls to close do not panic", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				c, _ := newCache(cacheSize)
				c.Close()
				require.NotPanics(t, func() { c.Close() })
			})
		}
	})

	t.Run("cache with no goroutines is garbage collected after releasing without closing", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				ch := make(chan struct{})

				cache := func() *memcache.Cache[int, int] {
					c, _ := newCache(cacheSize)
					runtime.SetFinalizer(c, func(_ *memcache.Cache[int, int]) {
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
		}
	})

	t.Run("cache with goroutines is garbage collected after releasing & closing", func(t *testing.T) {
		t.Parallel()

		for policy, newCache := range policies {
			newCache := newCache
			t.Run(policy, func(t *testing.T) {
				t.Parallel()

				ch := make(chan struct{})

				cache := func() *memcache.Cache[int, int] {
					interval := 1 * time.Second
					c, _ := newCache(cacheSize, memcache.WithActiveExpiration[int, int](interval))
					runtime.SetFinalizer(c, func(_ *memcache.Cache[int, int]) {
						close(ch)
					})
					return c
				}()

				cache.Flush() // use the cache once
				cache.Close() // close the cache
				cache = nil   // release the cache
				runtime.GC()  // explicitly run garbage collection

				select {
				case <-time.After(500 * time.Millisecond):
					t.Fatal("cache was not garbage collected")
				case <-ch:
					// cache was garbage collected
				}
			})
		}
	})
}

func TestCache_unsafe(t *testing.T) {
	t.Parallel()

	t.Run("demonstrates unsafe usage of pointer values stored in cache", func(t *testing.T) {
		t.Parallel()

		v := false
		c, _ := memcache.OpenNoEvictionCache[int, *bool]()

		c.Set(1, &v)
		v = true

		items, unlock := c.Store().Items()
		defer unlock()
		require.Contains(t, items, 1)
		require.Equal(t, true, *items[1].Value)
	})
}
