package memcache

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/wafer-bw/memcache/internal/closer"
)

type lruStore[K comparable, V any] struct {
	mu       *sync.RWMutex
	closer   *closer.Closer
	capacity int
	list     *list.List
	elements map[K]*list.Element
	items    map[K]Item[K, V]
	readerCh chan K
}

func newLRUStore[K comparable, V any](capacity int, closer *closer.Closer) (lruStore[K, V], error) {
	if capacity <= 1 {
		return lruStore[K, V]{}, ErrInvalidCapacity
	}

	store := lruStore[K, V]{
		mu:       &sync.RWMutex{},
		closer:   closer,
		capacity: capacity,
		list:     list.New(),
		elements: make(map[K]*list.Element, capacity),
		items:    make(map[K]Item[K, V], capacity),
		readerCh: make(chan K),
	}

	go store.reader()

	return store, nil
}

func (s lruStore[K, V]) Set(key K, value Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.elements) == s.capacity {
		element := s.list.Back()
		evictKey := element.Value.(K)

		// TODO: call item.OnEvicted:
		// item := s.items[key]
		// if item.OnEvicted != nil {
		// 	item.OnEvicted(key, item.Value)
		// }

		s.list.Remove(element)
		delete(s.elements, evictKey)
		delete(s.items, evictKey)

		fmt.Println(len(s.elements))
		fmt.Println(len(s.items))
		fmt.Println(s.list.Len())
	}

	element := s.list.PushFront(key)
	s.elements[key] = element
	s.items[key] = value
}

func (s lruStore[K, V]) Get(key K) (Item[K, V], bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[key]
	if !ok {
		return Item[K, V]{}, false
	}

	s.readerCh <- key

	return item, true
}

func (s lruStore[K, V]) Items() (map[K]Item[K, V], unlockFunc) {
	s.mu.Lock()
	return s.items, s.mu.Unlock
}

func (s lruStore[K, V]) Delete(keys ...K) {
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

func (s lruStore[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: call every item's OnEvicted:
	// item := s.items[key]
	// if item.OnEvicted != nil {
	// 	item.OnEvicted(key, item.Value)
	// }

	s.list.Init()
	clear(s.elements)
	clear(s.items)
}

func (s lruStore[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

func (s lruStore[K, V]) reader() {
	for {
		select {
		case <-s.closer.WaitClosed():
			return
		case key := <-s.readerCh:
			s.mu.Lock()
			s.list.MoveToFront(s.elements[key])
			s.mu.Unlock()
		}
	}
}
