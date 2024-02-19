package volatilelru

import (
	"container/list"
	"sync"

	"github.com/wafer-bw/memcache/internal/data"
)

const (
	PolicyName      string = "volatilelru"
	DefaultCapacity int    = 10_000
	MinimumCapacity int    = 2
)

type Store[K comparable, V any] struct {
	mu       sync.RWMutex
	capacity int
	items    map[K]data.Item[K, V]
	elements map[K]*list.Element
	list     *list.List
}

func New[K comparable, V any](capacity int) *Store[K, V] {
	if capacity < MinimumCapacity {
		capacity = DefaultCapacity
	}

	return &Store[K, V]{
		capacity: capacity,
		list:     list.New(),
		elements: make(map[K]*list.Element, capacity),
		items:    make(map[K]data.Item[K, V], capacity),
	}
}

func (s *Store[K, V]) Add(key K, item data.Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	element := s.list.PushFront(key)
	s.elements[key] = element
	s.items[key] = item

	if len(s.items) > s.capacity {
		s.evict()
	}
}

func (s *Store[K, V]) Get(key K) (data.Item[K, V], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return item, ok
	}

	s.list.MoveToFront(s.elements[key])

	return item, ok
}

func (s *Store[K, V]) Remove(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		s.delete(key)
	}
}

func (s *Store[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

func (s *Store[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	return keys
}

func (s *Store[K, V]) Items() map[K]data.Item[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make(map[K]data.Item[K, V], len(s.items))
	for key, item := range s.items {
		items[key] = item
	}

	return items
}

func (s *Store[K, V]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list.Init()
	clear(s.elements)
	clear(s.items)
}

func (s *Store[K, V]) evict() {
	cursor := s.list.Back()

	for {
		key, _ := cursor.Value.(K)
		item := s.items[key]

		if item.ExpireAt != nil {
			s.delete(key)
			return
		}

		cursor = cursor.Prev()
		if cursor == nil {
			break
		}
	}

	key, _ := s.list.Back().Value.(K)
	s.delete(key)
}

func (s *Store[K, V]) delete(key K) {
	element, ok := s.elements[key]
	if !ok {
		return
	}

	s.list.Remove(element)
	delete(s.elements, key)
	delete(s.items, key)
}
