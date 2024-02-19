// Package ports provides commonly shared internal interfaces.
package ports

import "github.com/wafer-bw/memcache/internal/data"

// TODO: docstring.
type Storer[K comparable, V any] interface {
	Add(key K, item data.Item[K, V])
	Get(key K) (data.Item[K, V], bool)
	Remove(keys ...K)
	Len() int
	Keys() []K
	Items() map[K]data.Item[K, V]
	Flush()
}

// TODO: docstring.
type Closer interface {
	Close()
	Closed() bool
}
