package noevict

import (
	"sync"
	"time"

	"github.com/wafer-bw/memcache/errs"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/store/closeable"
)

const (
	PolicyName      string = "noevict"
	MinimumCapacity int    = 0
)

// closer is the interface depended on by [Store] to ensure its goroutines are
// always closed.
type closer interface {
	Close()
	Closed() bool
}

type Store[K comparable, V any] struct {
	closer
	mu                sync.RWMutex
	capacity          int
	items             map[K]data.Item[K, V]
	passiveExpiration bool
}

type Config struct {
	Capacity                 int
	PassiveExpiration        bool
	ActiveExpirationInterval time.Duration
}

func Open[K comparable, V any](config Config) (*Store[K, V], error) {
	if config.Capacity < MinimumCapacity {
		return nil, errs.InvalidCapacityError{
			Policy:   PolicyName,
			Capacity: config.Capacity,
			Minimum:  MinimumCapacity,
		}
	}

	s := &Store[K, V]{
		closer:            closeable.New(),
		mu:                sync.RWMutex{},
		capacity:          config.Capacity,
		items:             make(map[K]data.Item[K, V], config.Capacity),
		passiveExpiration: config.PassiveExpiration,
	}

	if config.ActiveExpirationInterval > 0 {
		go s.runActiveExpirer(config.ActiveExpirationInterval)
	}

	return s, nil
}

func (s *Store[K, V]) Set(key K, value data.Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Do not add more keys when at capacity.
	if _, ok := s.items[key]; !ok && s.atCapacity() {
		return
	}

	s.items[key] = value
}

func (s *Store[K, V]) Get(key K) (data.Item[K, V], bool) {
	if s.passiveExpiration {
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
		if s.passiveExpiration {
			delete(s.items, key)
		}
		return data.Item[K, V]{}, false
	}

	return item, ok
}

func (s *Store[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.items, key)
	}
}

func (s *Store[K, V]) Size() int {
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

func (s *Store[K, V]) Items() (map[K]data.Item[K, V], func()) {
	s.mu.Lock()

	return s.items, s.mu.Unlock
}

func (s *Store[K, V]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.items)
}

func (s *Store[K, V]) expiredKeys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0, len(s.items))
	for key, item := range s.items {
		if item.IsExpired() {
			keys = append(keys, key)
		}
	}

	return keys
}

func (s *Store[K, V]) runActiveExpirer(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		if s.Closed() {
			return
		}
		s.Delete(s.expiredKeys()...)
	}
}

func (s *Store[K, V]) atCapacity() bool {
	return s.capacity > 0 && len(s.items) >= s.capacity
}
