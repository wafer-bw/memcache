package memcache_test

import (
	"fmt"
	"time"

	"github.com/wafer-bw/memcache"
)

func ExampleOpenNoEvictionCache() {
	capacity := 10
	interval := 1 * time.Second
	cache, err := memcache.OpenNoEvictionCache(
		memcache.WithActiveExpiration[int, string](interval),
		memcache.WithPassiveExpiration[int, string](),
		memcache.WithCapacity[int, string](capacity),
	)
	if err != nil {
		panic(err)
	}
	defer cache.Close()
}

func ExampleOpenAllKeysLRUCache() {
	capacity := 10
	interval := 1 * time.Second
	cache, err := memcache.OpenAllKeysLRUCache(capacity,
		memcache.WithActiveExpiration[int, string](interval),
		memcache.WithPassiveExpiration[int, string](),
	)
	if err != nil {
		panic(err)
	}
	defer cache.Close()
}

func ExampleOpenVolatileLRUCache() {
	capacity := 10
	interval := 1 * time.Second
	cache, err := memcache.OpenVolatileLRUCache(capacity,
		memcache.WithActiveExpiration[int, string](interval),
		memcache.WithPassiveExpiration[int, string](),
	)
	if err != nil {
		panic(err)
	}
	defer cache.Close()
}

func ExampleOpenAllKeysLFUCache() {
	capacity := 10
	interval := 1 * time.Second
	cache, err := memcache.OpenAllKeysLFUCache(capacity,
		memcache.WithActiveExpiration[int, string](interval),
		memcache.WithPassiveExpiration[int, string](),
	)
	if err != nil {
		panic(err)
	}
	defer cache.Close()
}

func ExampleCache_Set() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
}

func ExampleCache_SetEx() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.SetEx(1, "one", 1*time.Second)
}

func ExampleCache_Get() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")

	v, ok := cache.Get(1)
	fmt.Println(v)
	fmt.Println(ok)
	_, ok = cache.Get(2)
	fmt.Println(ok)
	// Output:
	// one
	// true
	// false
}

func ExampleCache_TTL() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.SetEx(1, "one", 2*time.Minute)

	v, ok := cache.TTL(1)
	fmt.Println(v.Truncate(time.Minute))
	fmt.Println(ok)
	// Output:
	// 1m0s
	// true
}

func ExampleCache_Delete() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}
	cache.Set(1, "one")
	cache.Set(2, "two")
	cache.Set(3, "three")

	cache.Delete(1)
	cache.Delete(2, 3)

	_, ok := cache.Get(1)
	fmt.Println(ok)
	_, ok = cache.Get(2)
	fmt.Println(ok)
	_, ok = cache.Get(3)
	fmt.Println(ok)
	// Output:
	// false
	// false
	// false
}

func ExampleCache_Flush() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
	cache.Set(2, "two")

	cache.Flush()

	_, ok := cache.Get(1)
	fmt.Println(ok)

	_, ok = cache.Get(2)
	fmt.Println(ok)
	// Output:
	// false
	// false
}

func ExampleCache_Size() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
	cache.Set(2, "two")

	fmt.Println(cache.Size())
	// Output:
	// 2
}

func ExampleCache_RandomKey() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
	cache.Set(2, "two")

	key, ok := cache.RandomKey()
	_ = key
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleCache_Keys() {
	cache, err := memcache.OpenNoEvictionCache[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
	cache.Set(2, "two")

	for _, key := range cache.Keys() {
		fmt.Println(key)
	}
	// Unordered output:
	// 1
	// 2
}
