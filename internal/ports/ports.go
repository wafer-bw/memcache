// Package ports provides commonly shared internal interfaces.
package ports

import (
	"time"

	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/expire"
)

type Storer[K comparable, V any] interface {
	Add(key K, item data.Item[K, V])
	Get(key K) (data.Item[K, V], bool)
	Remove(keys ...K)
	Len() int
	Keys() []K
	Items() map[K]data.Item[K, V]
	Flush()
}

type Cacher[K comparable, V any] interface {
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
	// - Random()  // return random key/value from cache.
	// Maybe:
	// - Persist() // remove ttl from key.
	// - Expire()  // set ttl for key.
}

type Closer interface {
	Close()
	Closed() bool
	Ch() <-chan struct{}
}

type Expirer[K comparable, V any] interface {
	Expire(expire.Cacher[K, V])
}
