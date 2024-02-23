package allkeyslfu

import (
	"sync"

	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/substore/lfulist"
	"github.com/wafer-bw/memcache/internal/substore/randxs"
)

const (
	PolicyName      string = "allkeyslru"
	DefaultCapacity int    = 10_000
	MinimumCapacity int    = 2
)

type Store[K comparable, V any] struct {
	mu       sync.RWMutex
	capacity int

	items        map[K]data.Item[K, V]   // primary storage of key-value pairs
	randomAccess ports.RandomAccessor[K] // permits random key selection
	lfu          ports.LFUTracker[K]     // permits least frequently used key selection
}

func New[K comparable, V any](capacity int) *Store[K, V] {
	if capacity < MinimumCapacity {
		capacity = DefaultCapacity
	}

	return &Store[K, V]{
		capacity:     capacity,
		items:        make(map[K]data.Item[K, V], capacity),
		randomAccess: randxs.New[K](capacity),
		lfu:          lfulist.New[K](capacity),
	}
}

func (s *Store[K, V]) Add(key K, item data.Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.randomAccess.Add(key)
	s.items[key] = item
	s.lfu.Inc(key)

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

	s.lfu.Inc(key)

	return item, ok
}

func (s *Store[K, V]) Remove(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		s.delete(key)
		s.lfu.Remove(key)
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
	s.randomAccess.Clear()
	s.lfu.Clear()
}

func (s *Store[K, V]) evict() {
	s.delete(s.lfu.LFU())
}

func (s *Store[K, V]) delete(key K) {
	delete(s.items, key)
	s.randomAccess.Remove(key)
	s.lfu.Remove(key)
}
