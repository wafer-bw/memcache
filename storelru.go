package memcache

import (
	"container/list"
	"sync"
)

type lruStore[K comparable, V any] struct {
	mu sync.RWMutex

	capacity int
	list     *list.List
	elements map[K]*list.Element
	items    map[K]Item[K, V]
}

func newLRUStore[K comparable, V any](capacity int) (*lruStore[K, V], error) {
	if capacity <= 1 {
		return nil, ErrInvalidCapacity
	}

	store := &lruStore[K, V]{
		capacity: capacity,
		list:     list.New(),
		elements: make(map[K]*list.Element, capacity),
		items:    make(map[K]Item[K, V], capacity),
	}

	return store, nil
}

func (s *lruStore[K, V]) Set(key K, value Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	element := s.list.PushFront(key)
	s.elements[key] = element
	s.items[key] = value

	if len(s.elements) > s.capacity {
		element := s.list.Back()
		evictKey := element.Value.(K)

		s.list.Remove(element)
		delete(s.elements, evictKey)
		delete(s.items, evictKey)
	}
}

func (s *lruStore[K, V]) Get(key K, deleteExpired bool) (Item[K, V], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return Item[K, V]{}, false
	}

	if item.IsExpired() && deleteExpired {
		s.list.Remove(s.elements[key])
		delete(s.elements, key)
		delete(s.items, key)

		return Item[K, V]{}, false
	}

	s.list.MoveToFront(s.elements[key])

	return item, true
}

func (s *lruStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.mu.Lock()
	return s.items, s.mu.Unlock
}

func (s *lruStore[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	return keys
}

func (s *lruStore[K, V]) Delete(keys ...K) {
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

func (s *lruStore[K, V]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list.Init()
	clear(s.elements)
	clear(s.items)
}

func (s *lruStore[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}
