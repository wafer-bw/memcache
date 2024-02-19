package memcache_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var sizes = []int{100, 1000, 10000, 100000}

func BenchmarkCache_Set(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					cache.Set(i%size, i%size)
				}
			})
		}
	}
}

func BenchmarkCache_Get(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = cache.Get(i % size)
				}
			})
		}
	}
}

func BenchmarkCache_TTL(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.SetEx(i, i, time.Duration(rand.Intn(100))*time.Millisecond)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = cache.TTL(i % size)
				}
			})
		}
	}
}

func BenchmarkCache_Delete(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					cache.Delete(i % size)
				}
			})
		}
	}
}

func BenchmarkCache_Flush(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					cache.Flush()
				}
			})
		}
	}
}

func BenchmarkCache_Size(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = cache.Size()
				}
			})
		}
	}
}

func BenchmarkCache_Keys(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d %s", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = cache.Keys()
				}
			})
		}
	}
}
