package memcache_test

import (
	"fmt"
	"testing"
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

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
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

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = cache.Get(i % size)
				}
			})
		}
	}
}

func BenchmarkCache_Get_parallel(b *testing.B) {
	for policy, newCache := range policies {
		cache, err := newCache(2)
		if err != nil {
			b.Fatal(err)
		}

		cache.Set(1, 1)

		b.Run(fmt.Sprintf("%s policy parallel", policy), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = cache.Get(1)
				}
			})
		})
	}
}

func BenchmarkCache_Has(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = cache.Has(i % size)
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

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
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

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
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

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = cache.Size()
				}
			})
		}
	}
}

func BenchmarkCache_Items(b *testing.B) {
	for policy, newCache := range policies {
		for _, size := range sizes {
			cache, err := newCache(size)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < size; i++ {
				cache.Set(i, i)
			}

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, unlock := cache.Store().Items()
					unlock()
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

			b.Run(fmt.Sprintf("%d items %s policy", size, policy), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = cache.Keys()
				}
			})
		}
	}
}
