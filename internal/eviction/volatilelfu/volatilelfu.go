package volatilelfu

import (
	"sync"

	"github.com/wafer-bw/memcache/internal/data"
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/substore/randxs"
)

const (
	PolicyName      string = "volatilelfu"
	DefaultCapacity int    = 10_000
	MinimumCapacity int    = 2
)

type Store[K comparable, V any] struct {
	mu       sync.RWMutex
	capacity int

	items        map[K]data.Item[K, V]   // primary storage of key-value pairs
	randomAccess ports.RandomAccessor[K] // permits random key selection
	nodes        map[K]*freqNode[K]      // lfu frequency structure
	frequencies  map[int]*list[K]        // lfu frequency structure
	min          int                     // lfu frequency structure
}

func New[K comparable, V any](capacity int) *Store[K, V] {
	if capacity < MinimumCapacity {
		capacity = DefaultCapacity
	}

	return &Store[K, V]{
		capacity:     capacity,
		items:        make(map[K]data.Item[K, V], capacity),
		randomAccess: randxs.New[K](capacity),
		nodes:        make(map[K]*freqNode[K], capacity),
		frequencies:  make(map[int]*list[K], capacity),
		min:          0,
	}
}

func (s *Store[K, V]) Add(key K, item data.Item[K, V]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.randomAccess.Add(key)
	s.items[key] = item
	s.inc(key)

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

	s.inc(key)

	return item, ok
}

func (s *Store[K, V]) Remove(keys ...K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		s.delete(key)
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
	clear(s.nodes)
	clear(s.frequencies)
	s.min = 0
}

func (s *Store[K, V]) inc(key K) {
	node, ok := s.nodes[key]
	if !ok {
		node := &freqNode[K]{key: key, freq: 1}
		s.min = 1
		list, ok := s.frequencies[node.freq]
		if !ok {
			list = newList[K]()
		}
		list.pushBack(node)
		s.frequencies[node.freq] = list
		s.nodes[key] = node
		return
	}
	list, ok := s.frequencies[node.freq]
	if ok {
		list.remove(node)
	}
	node.freq++

	nextList, ok := s.frequencies[node.freq]
	if !ok {
		nextList = newList[K]()
	}

	nextList.pushBack(node)
	s.frequencies[node.freq] = nextList

	if list.size == 0 && s.min == node.freq-1 {
		s.min++
	}
}

func (s *Store[K, V]) evict() {
	s.delete(s.frequencies[s.min].head.next.key) // TODO: need to evict least frequently used volatile key
}

func (s *Store[K, V]) delete(key K) {
	delete(s.items, key)
	s.randomAccess.Remove(key)
	if node, ok := s.nodes[key]; ok {
		s.frequencies[node.freq].remove(node)
	}
	delete(s.nodes, key)
}

type list[K comparable] struct {
	// TODO: replace with go linked list using list.List & list.Elements,
	//       will need to be its own substore package.

	head *freqNode[K]
	tail *freqNode[K]
	size int
}

type freqNode[K comparable] struct {
	prev, next *freqNode[K]
	key        K
	freq       int
}

func newList[K comparable]() *list[K] {
	head := &freqNode[K]{}
	tail := &freqNode[K]{}
	head.next = tail
	tail.prev = head
	return &list[K]{
		head: head,
		tail: tail,
	}
}

func (l *list[K]) pushBack(node *freqNode[K]) {
	node.prev = l.tail.prev
	node.next = l.tail
	l.tail.prev.next = node
	l.tail.prev = node
	l.size++
}

func (l *list[K]) remove(node *freqNode[K]) {
	node.prev.next = node.next
	node.next.prev = node.prev
	l.size--
}
