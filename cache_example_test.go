package memcache_test

import (
	"fmt"
	"time"

	"github.com/wafer-bw/memcache"
)

// Create a new cache using int keys and string values.
func ExampleNew() {
	cache, err := memcache.Open[int, string]()
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleNew_withPassiveExpirationEnabled() {
	cache, err := memcache.Open[int, string](memcache.WithPassiveExpiration[int, string]())
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleNew_withExpirer() {
	interval := 1 * time.Second
	expirer := memcache.DeleteAllExpired[int, string]
	cache, err := memcache.Open[int, string](memcache.WithExpirer[int, string](expirer, interval))
	if err != nil {
		panic(err)
	}
	_ = cache
}

func ExampleCache_Set() {
	cache, err := memcache.Open[int, string]()
	if err != nil {
		panic(err)
	}

	cache.Set(1, "one")
}

func ExampleCache_SetEx() {
	cache, err := memcache.Open[int, string]()
	if err != nil {
		panic(err)
	}

	cache.SetEx(1, "one", 1*time.Second)
}

func ExampleCache_Get() {
	cache, err := memcache.Open[int, string]()
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
	cache, err := memcache.Open[int, string]()
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
	cache, err := memcache.Open[int, string]()
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
	cache, err := memcache.Open[int, string]()
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
	cache, err := memcache.Open[int, string]()
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
	cache, err := memcache.Open[int, string]()
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
