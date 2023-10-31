package memcache

// Everything in this file is exported for testing purposes only.

import (
	"container/list"
)

type UnlockFunc unlockFunc

func (c *Cache[K, V]) Store() storer[K, V] {
	return c.store
}

func (c *Cache[K, V]) PassiveExpiration() bool {
	return c.passiveExpiration
}

func (c *Cache[K, V]) Closed() bool {
	return c.closer.Closed()
}

func (c *Cache[K, V]) Expirer() ExpirerFunc[K, V] {
	return c.expirer
}

type LRUStore[K comparable, V any] struct {
	Underlying lruStore[K, V]
}

func NewLRUStore[K comparable, V any](capacity int) (LRUStore[K, V], error) {
	store, err := newLRUStore[K, V](capacity)
	if err != nil {
		return LRUStore[K, V]{}, err
	}

	return LRUStore[K, V]{
		Underlying: store,
	}, nil
}

func (s LRUStore[K, V]) Items() map[K]Item[K, V] {
	return s.Underlying.items
}

func (s LRUStore[K, V]) Keys() []K {
	return s.Underlying.Keys()
}

func (s LRUStore[K, V]) Elements() map[K]*list.Element {
	return s.Underlying.elements
}

func (s LRUStore[K, V]) List() *list.List {
	return s.Underlying.list
}

type NoEvictStore[K comparable, V any] struct {
	Underlying noEvictStore[K, V]
}

func NewNoEvictStore[K comparable, V any]() NoEvictStore[K, V] {
	return NoEvictStore[K, V]{
		Underlying: newNoEvictStore[K, V](),
	}
}

func (s NoEvictStore[K, V]) Items() map[K]Item[K, V] {
	return s.Underlying.items
}

func (s NoEvictStore[K, V]) Keys() []K {
	return s.Underlying.Keys()
}
