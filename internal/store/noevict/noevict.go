package noevict

import (
	"sync"

	"github.com/wafer-bw/memcache/internal/data"
)

type Store[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]data.Item[K, V]
}

func Open[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		items: make(map[K]data.Item[K, V]),
	}
}

func (s *Store[K, V]) Set(key K, value data.Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
}

func (s *Store[K, V]) Get(key K, deleteExpired bool) (data.Item[K, V], bool) {
	if deleteExpired {
		s.mu.Lock()
		defer s.mu.Unlock()
	} else {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	item, ok := s.items[key]
	if !ok {
		return data.Item[K, V]{}, false
	}

	if item.IsExpired() {
		if deleteExpired {
			delete(s.items, key)
		}
		return data.Item[K, V]{}, false
	}

	return item, ok
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

func (s *Store[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.items, key)
	}
}

func (s *Store[K, V]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.items)
}

func (s *Store[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

func (s *Store[K, V]) Items() (map[K]data.Item[K, V], func()) {
	s.mu.Lock()

	return s.items, s.mu.Unlock
}
