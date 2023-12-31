package memcache

import "sync"

type noEvictStore[K comparable, V any] struct {
	mu sync.RWMutex

	items map[K]Item[K, V]
}

func newNoEvictStore[K comparable, V any]() *noEvictStore[K, V] {
	return &noEvictStore[K, V]{
		items: make(map[K]Item[K, V]),
	}
}

func (s *noEvictStore[K, V]) Set(key K, value Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
}

func (s *noEvictStore[K, V]) Get(key K, deleteExpired bool) (Item[K, V], bool) {
	if deleteExpired {
		s.mu.Lock()
		defer s.mu.Unlock()
	} else {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	item, ok := s.items[key]
	if !ok {
		return Item[K, V]{}, false
	}

	if item.IsExpired() && deleteExpired {
		delete(s.items, key)
		return Item[K, V]{}, false
	}

	return item, ok
}

func (s *noEvictStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.mu.Lock()

	return s.items, s.mu.Unlock
}

func (s *noEvictStore[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	return keys
}

func (s *noEvictStore[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.items, key)
	}
}

func (s *noEvictStore[K, V]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.items)
}

func (s *noEvictStore[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}
