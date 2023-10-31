package memcache

import "sync"

type noEvictStore[K comparable, V any] struct {
	mu    *sync.RWMutex
	items map[K]Item[K, V]
}

func newNoEvictStore[K comparable, V any]() noEvictStore[K, V] {
	return noEvictStore[K, V]{
		mu:    &sync.RWMutex{},
		items: make(map[K]Item[K, V]),
	}
}

func (s noEvictStore[K, V]) Set(key K, value Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
}

func (s noEvictStore[K, V]) Get(key K, activelyExpire bool) (Item[K, V], bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[key]
	if !ok {
		return Item[K, V]{}, false
	}

	if activelyExpire && item.IsExpired() {
		return Item[K, V]{}, false
	}

	return item, ok
}

func (s noEvictStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.mu.Lock()

	return s.items, s.mu.Unlock
}

func (s noEvictStore[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.items, key)
	}
}

func (s noEvictStore[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.items)
}

func (s noEvictStore[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}
