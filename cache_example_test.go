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

func ExampleNew_withExpirationInterval() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx, memcache.WithExpirationInterval(1*time.Minute))
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

	cache.Set(1, "one", memcache.WithTTL(1*time.Second))
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
	// Output:
	// one
	// true
}

func ExampleCache_Get_keyNotFound() {
	ctx := context.TODO()

	cache, err := memcache.New[int, string](ctx)
	if err != nil {
		panic(err)
	}

	_, ok := cache.Get(1)
	fmt.Println(ok)
	// Output:
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

func ExampleCache_Length() {
	// TODO: this
}

func ExampleCache_Keys() {
	// TODO: this
}
