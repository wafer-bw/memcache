package lfulist

// TODO: this may be better just as part of an eviction policy store because
//       node contents may change when we need to include ttls

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

type Store[K comparable] struct {
	nodes       map[K]*freqNode[K]
	frequencies map[int]*list[K]
	capacity    int
	min         int
}

func New[K comparable](capacity int) *Store[K] {
	return &Store[K]{
		nodes:       make(map[K]*freqNode[K], capacity),
		frequencies: make(map[int]*list[K], capacity),
		capacity:    capacity,
		min:         0,
	}
}

func (s *Store[K]) Inc(key K) {
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

func (s *Store[K]) Remove(key K) {
	node, ok := s.nodes[key]
	if !ok {
		return
	}

	s.frequencies[node.freq].remove(node)
	delete(s.nodes, key)
}

func (s *Store[K]) LFU() K {
	return s.frequencies[s.min].head.next.key
}

func (s *Store[K]) Clear() {
	clear(s.nodes)
	clear(s.frequencies)
	s.min = 0
}
