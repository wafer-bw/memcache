package memcache

// Everything in this file is exported for testing purposes only.

import "container/list"

// UnlockFunc unlockes the mutex for the cache store.
type UnlockFunc func()

func (c *Cache[K, V]) Store() (storer[K, V], UnlockFunc) {
	c.mu.Lock()

	return c.store, c.mu.Unlock
}

func (c *Cache[K, V]) Items() (map[K]Item[K, V], UnlockFunc) {
	c.mu.Lock()

	return c.store.Items(), c.mu.Unlock
}

func (c *Cache[K, V]) Lock() {
	c.mu.Lock()
}

func (c *Cache[K, V]) Unlock() {
	c.mu.Unlock()
}

func (c *Cache[K, V]) RLock() {
	c.mu.RLock()
}

func (c *Cache[K, V]) RUnlock() {
	c.mu.RUnlock()
}

func (c *Cache[K, V]) PassiveExpiration() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.passiveExpiration
}

func (c *Cache[K, V]) Closed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.closed
}

func (c *Cache[K, V]) Expirer() ExpirerFunc[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()

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
