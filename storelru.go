package memcache

import (
	"container/list"
	"sync"
)

type lruStore[K comparable, V any] struct {
	mu       *sync.RWMutex
	capacity int
	list     *list.List
	elements map[K]*list.Element
	items    map[K]Item[K, V]
	readerCh chan K
}

func newLRUStore[K comparable, V any](capacity int) (lruStore[K, V], error) {
	if capacity <= 1 {
		return lruStore[K, V]{}, ErrInvalidCapacity
	}

	store := lruStore[K, V]{
		mu:       &sync.RWMutex{},
		capacity: capacity,
		list:     list.New(),
		elements: make(map[K]*list.Element, capacity),
		items:    make(map[K]Item[K, V], capacity),
		readerCh: make(chan K),
	}

	return store, nil
}

func (s lruStore[K, V]) Set(key K, value Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	element := s.list.PushFront(key)
	s.elements[key] = element
	s.items[key] = value

	if len(s.elements) > s.capacity {
		element := s.list.Back()
		evictKey := element.Value.(K)

		// TODO: call item.OnEvicted:
		// item := s.items[key]
		// if item.OnEvicted != nil {
		// 	item.OnEvicted(key, item.Value)
		// }

		s.list.Remove(element)
		delete(s.elements, evictKey)
		delete(s.items, evictKey)
	}
}

func (s lruStore[K, V]) Get(key K, activelyExpire bool) (Item[K, V], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return Item[K, V]{}, false
	}

	if activelyExpire && item.IsExpired() {
		s.list.Remove(s.elements[key])
		delete(s.elements, key)
		delete(s.items, key)

		return Item[K, V]{}, false
	}

	s.list.MoveToFront(s.elements[key])

	return item, true
}

func (s lruStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.mu.Lock()
	return s.items, s.mu.Unlock
}

func (s lruStore[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		element, ok := s.elements[key]
		if !ok {
			continue
		}

		s.list.Remove(element)
		delete(s.elements, key)
		delete(s.items, key)
	}
}

func (s lruStore[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: call every item's OnEvicted:
	// item := s.items[key]
	// if item.OnEvicted != nil {
	// 	item.OnEvicted(key, item.Value)
	// }

	s.list.Init()
	clear(s.elements)
	clear(s.items)
}

func (s lruStore[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}
