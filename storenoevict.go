package memcache

import "sync"

type noEvictStore[K comparable, V any] struct {
	mu    *sync.RWMutex
	items map[K]Item[K, V]
	keys  map[K]struct{}
}

func newNoEvictStore[K comparable, V any]() noEvictStore[K, V] {
	return noEvictStore[K, V]{
		mu:    &sync.RWMutex{},
		items: make(map[K]Item[K, V]),
		keys:  make(map[K]struct{}),
	}
}

func (s noEvictStore[K, V]) Set(key K, value Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
	s.keys[key] = struct{}{}
}

func (s noEvictStore[K, V]) Get(key K, passivelyExpire bool) (Item[K, V], bool) {
	if passivelyExpire {
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

	if item.IsExpired() {
		if passivelyExpire {
			delete(s.items, key)
			delete(s.keys, key)
		}
		return Item[K, V]{}, false
	}

	return item, ok
}

func (s noEvictStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.mu.Lock()

	return s.items, s.mu.Unlock
}

func (s noEvictStore[K, V]) Keys() map[K]struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.keys
}

func (s noEvictStore[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.items, key)
		delete(s.keys, key)
	}
}

func (s noEvictStore[K, V]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.items)
	clear(s.keys)
}

func (s noEvictStore[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}
