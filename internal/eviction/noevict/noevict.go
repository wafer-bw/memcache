package noevict

import (
	"sync"

	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/substore/randxs"
)

const (
	PolicyName      string = "noevict"
	DefaultCapacity int    = 0
	MinimumCapacity int    = 0
)

type Store[K comparable, V any] struct {
	mu           sync.RWMutex
	capacity     int
	items        map[K]data.Item[K, V]   // primary storage of key-value pairs
	randomAccess ports.RandomAccessor[K] // permits random key selection
}

func New[K comparable, V any](capacity int) *Store[K, V] {
	if capacity < 0 {
		capacity = DefaultCapacity
	}

	return &Store[K, V]{
		mu:           sync.RWMutex{},
		capacity:     capacity,
		items:        make(map[K]data.Item[K, V], capacity),
		randomAccess: randxs.New[K](capacity),
	}
}

func (s *Store[K, V]) Add(key K, item data.Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.atCapacity() {
		return
	}

	s.randomAccess.Add(key)
	s.items[key] = item
}

func (s *Store[K, V]) Get(key K) (data.Item[K, V], bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[key]
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

func (s *Store[K, V]) RandomKey() (K, bool) {
	return s.randomAccess.RandomKey()
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

	clear(s.items)
}

func (s *Store[K, V]) delete(key K) {
	s.randomAccess.Remove(key)
	delete(s.items, key)
}

func (s *Store[K, V]) atCapacity() bool {
	return s.capacity > 0 && len(s.items) >= s.capacity
}
