package lru

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/wafer-bw/memcache/internal/closeable"
	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/expire"
)

const (
	PolicyName      string = "lru"
	MinimumCapacity int    = 1
)

type closer interface {
	Close()
	Closed() bool
}

type expirer[K comparable, V any] interface {
	Expire(expire.Storer[K, V])
}

type Store[K comparable, V any] struct {
	closer
	expirer           expirer[K, V]
	mu                sync.RWMutex
	keys              []K
	list              *list.List
	elements          map[K]*list.Element
	items             map[K]data.Item[K, V]
	capacity          int
	passiveExpiration bool
}

type Config struct {
	PassiveExpiration        bool
	ActiveExpirationInterval time.Duration
}

func Open[K comparable, V any](capacity int, config Config) (*Store[K, V], error) {
	if capacity < MinimumCapacity {
		return nil, errors.New("invalid capacity")
	}

	s := &Store[K, V]{
		closer:            closeable.New(),
		expirer:           expire.AllKeys[K, V]{},
		mu:                sync.RWMutex{},
		capacity:          capacity,
		keys:              make([]K, 0, capacity),
		list:              list.New(),
		elements:          make(map[K]*list.Element, capacity),
		items:             make(map[K]data.Item[K, V], capacity),
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

	// Adding a new key at capacity requires eviction.
	if _, ok := s.items[key]; !ok && s.atCapacity() {
		element := s.list.Back()
		evictKey, _ := element.Value.(K)

		s.list.Remove(element)
		delete(s.elements, evictKey)
		delete(s.items, evictKey)
	}

	element := s.list.PushFront(key)
	s.elements[key] = element
	s.items[key] = value
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

func (s *Store[K, V]) runActiveExpirer(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		if s.Closed() {
			return
		}
		s.expirer.Expire(s)
	}
}

func (s *Store[K, V]) atCapacity() bool {
	return len(s.items) >= s.capacity
}
