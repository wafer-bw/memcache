package memcache

import "sync"

type noEvictStore[K comparable, V any] struct {
	*sync.RWMutex

	items map[K]Item[K, V]
	keys  map[K]struct{}
}

func newNoEvictStore[K comparable, V any]() noEvictStore[K, V] {
	return noEvictStore[K, V]{
		RWMutex: &sync.RWMutex{},
		items:   make(map[K]Item[K, V]),
		keys:    make(map[K]struct{}),
	}
}

func (s noEvictStore[K, V]) Set(key K, value Item[K, V]) {
	s.Lock()
	defer s.Unlock()

	s.items[key] = value
	s.keys[key] = struct{}{}
}

func (s noEvictStore[K, V]) Get(key K, deleteExpired bool) (Item[K, V], bool) {
	if deleteExpired {
		s.Lock()
		defer s.Unlock()
	} else {
		s.RLock()
		defer s.RUnlock()
	}

	item, ok := s.items[key]
	if !ok {
		return Item[K, V]{}, false
	}

	if item.IsExpired() {
		if deleteExpired {
			delete(s.items, key)
			delete(s.keys, key)
		}
		return Item[K, V]{}, false
	}

	return item, ok
}

func (s noEvictStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.Lock()

	return s.items, s.Unlock
}

func (s noEvictStore[K, V]) Keys() map[K]struct{} {
	s.RLock()
	defer s.RUnlock()

	return s.keys
}

func (s noEvictStore[K, V]) Delete(keys ...K) {
	s.Lock()
	defer s.Unlock()

	for _, key := range keys {
		delete(s.items, key)
		delete(s.keys, key)
	}
}

func (s noEvictStore[K, V]) Flush() {
	s.Lock()
	defer s.Unlock()

	clear(s.items)
	clear(s.keys)
}

func (s noEvictStore[K, V]) Size() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.items)
}
