package memcache_test

import (
	"fmt"
	"time"

	"github.com/wafer-bw/memcache"
)

// Create a new cache using int keys and string values.
func ExampleNew() {
	cache, err := memcache.New[int, string]()
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleNew_withExpirationInterval() {
	cache, err := memcache.New[int, string](memcache.WithExpirationInterval(1 * time.Minute))
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleCache_Set() {
	cache, err := memcache.New[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
}

func ExampleCache_Set_withTTL() {
	cache, err := memcache.New[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one", memcache.WithTTL(1*time.Second))
}

func ExampleCache_Get() {
	cache, err := memcache.New[int, string]()
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
	cache, err := memcache.New[int, string]()
	if err != nil {
		panic(err)
	}

	_, ok := cache.Get(1)
	fmt.Println(ok)
	// Output:
	// false
}

func ExampleCache_Delete() {
	cache, err := memcache.New[int, string]()
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
	cache, err := memcache.New[int, string]()
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
