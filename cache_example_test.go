package memcache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/wafer-bw/memcache"
)

// Create a new cache using int keys and string values.
func ExampleNew() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleNew_withPassiveExpirationEnabled() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx, memcache.WithPassiveExpiration[int, string]())
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleCache_Set() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
}

func ExampleCache_Set_withTTL() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one", memcache.WithTTL[int, string](1*time.Second))
}

func ExampleCache_Get() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
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

func ExampleCache_Has() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")

	fmt.Println(cache.Has(1))
	fmt.Println(cache.Has(2))
	// Output:
	// true
	// false
}

func ExampleCache_Delete() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")

	cache.Delete(1)

	_, ok := cache.Get(1)
	fmt.Println(ok)
	// Output:
	// false
}

func ExampleCache_Flush() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
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
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
	cache.Set(2, "two")

	fmt.Println(cache.Size())
	// Output:
	// 2
}

func ExampleCache_Keys() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
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
