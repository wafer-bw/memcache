package memcache

type noEvictStore[K comparable, V any] struct {
	items map[K]Item[K, V]
}

func newNoEvictStore[K comparable, V any]() noEvictStore[K, V] {
	return noEvictStore[K, V]{items: make(map[K]Item[K, V])}
}

func (s noEvictStore[K, V]) Set(key K, value Item[K, V]) {
	s.items[key] = value
}

func (s noEvictStore[K, V]) Get(key K) (Item[K, V], bool) {
	v, ok := s.items[key]
	return v, ok
}

func (s noEvictStore[K, V]) Items() map[K]Item[K, V] {
	return s.items
}

func (s noEvictStore[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(s.items, key)
	}
}

func (s noEvictStore[K, V]) Clear() {
	clear(s.items)
}

func (s noEvictStore[K, V]) Size() int {
	return len(s.items)
}
