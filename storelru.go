package memcache

import "container/list"

type lruStore[K comparable, V any] struct {
	capacity int
	list     *list.List
	elements map[K]*list.Element
	items    map[K]Item[K, V]
}

func newLRUStore[K comparable, V any](capacity int) (lruStore[K, V], error) {
	if capacity <= 1 {
		return lruStore[K, V]{}, ErrInvalidCapacity
	}

	return lruStore[K, V]{
		capacity: capacity,
		list:     list.New(),
		elements: make(map[K]*list.Element, capacity),
		items:    make(map[K]Item[K, V], capacity),
	}, nil
}

func (s lruStore[K, V]) Set(key K, value Item[K, V]) {
	s.eviction()
	element := s.list.PushFront(key)
	s.elements[key] = element
	s.items[key] = value
}

func (s lruStore[K, V]) Get(key K) (Item[K, V], bool) {
	element, ok := s.elements[key]
	if !ok {
		return Item[K, V]{}, false
	}
	s.list.MoveToFront(element)
	key = element.Value.(K)
	item := s.items[key] // TODO: check !ok and remove from elements & list
	return item, true
}

func (s lruStore[K, V]) Items() map[K]Item[K, V] {
	return s.items
}

func (s lruStore[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		element := s.elements[key] // TODO: check !ok and handle?

		s.list.Remove(element)
		delete(s.elements, key)
		delete(s.items, key)
	}
}

func (s lruStore[K, V]) Clear() {
	s.list.Init()
	clear(s.elements)
	clear(s.items)
}

func (s lruStore[K, V]) Size() int {
	return len(s.items)
}

func (s lruStore[K, V]) eviction() {
	if len(s.elements) == s.capacity {
		element := s.list.Back()
		key := element.Value.(K)

		// TODO: add this.
		// if storeItem.value.OnEvicted != nil {
		// 	storeItem.value.OnEvicted(storeItem.key, storeItem.value.Value)
		// }

		s.list.Remove(element)
		delete(s.elements, key)
		delete(s.items, key)
	}
}
