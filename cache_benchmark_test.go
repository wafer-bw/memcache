package memcache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/wafer-bw/memcache"
)

func BenchmarkCache_Set(b *testing.B) {
	ctx := context.Background()

	for _, n := range []int{100, 1000, 10000, 100000} {
		cache, err := memcache.New[int, int](ctx)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < n; i++ {
			cache.Set(i, i)
		}

		b.Run(fmt.Sprintf("set to cache with %d keys", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Set(i%n, i%n)
			}
		})
	}
}

func BenchmarkCache_Get(b *testing.B) {
	ctx := context.Background()

	for _, n := range []int{100, 1000, 10000, 100000} {
		cache, err := memcache.New[int, int](ctx)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < n; i++ {
			cache.Set(i, i)
		}

		b.Run(fmt.Sprintf("get from cache with %d keys", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = cache.Get(i % n)
			}
		})
	}
}

func BenchmarkCache_Delete(b *testing.B) {
	ctx := context.Background()

	for _, n := range []int{100, 1000, 10000, 100000} {
		cache, err := memcache.New[int, int](ctx)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < n; i++ {
			cache.Set(i, i)
		}

		b.Run(fmt.Sprintf("delete from cache with %d keys", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Delete(i % n)
			}
		})
	}
}

func BenchmarkCache_Flush(b *testing.B) {
	ctx := context.Background()

	for _, n := range []int{100, 1000, 10000, 100000} {
		cache, err := memcache.New[int, int](ctx)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < n; i++ {
			cache.Set(i, i)
		}

		b.Run(fmt.Sprintf("flush cache with %d keys", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Flush()
			}
		})
	}
}
