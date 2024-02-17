package lru

import (
	"container/list"
	"sync"
	"time"

	"github.com/wafer-bw/memcache/internal/data"
)

const (
	PolicyName      string = "lru"
	MinimumCapacity int    = 1
)

type Store[K comparable, V any] struct {
	mu                sync.RWMutex
	closeCh           chan struct{}
	list              *list.List
	elements          map[K]*list.Element
	items             map[K]data.Item[K, V]
	capacity          int
	passiveExpiration bool
}

type Config struct {
	Capacity                 int
	PassiveExpiration        bool
	ActiveExpirationInterval time.Duration
}

func Open[K comparable, V any](config Config) (*Store[K, V], error) {
	s := &Store[K, V]{
		mu:                sync.RWMutex{},
		closeCh:           make(chan struct{}),
		capacity:          config.Capacity,
		list:              list.New(),
		elements:          make(map[K]*list.Element, config.Capacity),
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

func (s *Store[K, V]) Get(key K) (data.Item[K, V], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return data.Item[K, V]{}, false
	}

	if item.IsExpired() {
		if s.passiveExpiration {
			s.list.Remove(s.elements[key])
			delete(s.elements, key)
			delete(s.items, key)
		}
		return data.Item[K, V]{}, false
	}

	s.list.MoveToFront(s.elements[key])

	return item, true
}

func (s *Store[K, V]) Delete(keys ...K) {
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

	s.list.Init()
	clear(s.elements)
	clear(s.items)
}

func (s *Store[K, V]) Close() {
	// TODO: does this need to be protected by a mutex?
	select {
	case <-s.closeCh:
		return
	default:
		close(s.closeCh)
	}
}

func (s *Store[K, V]) Closed() bool {
	// TODO: does this need to be protected by a mutex?
	select {
	case <-s.closeCh:
		return true
	default:
		return false
	}
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
		select {
		case <-s.closeCh:
			return
		case <-ticker.C:
			s.Delete(s.expiredKeys()...)
		}
	}
}
