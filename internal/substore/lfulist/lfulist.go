package lfulist

// TODO: replace with go linked list using list.List & list.Elements,
// will need to be its own substore package.
type list[K comparable] struct {
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
	nodeMap  map[K]*freqNode[K]
	listMap  map[int]*list[K]
	capacity int
	min      int
}

func New[K comparable](capacity int) *Store[K] {
	return &Store[K]{
		nodeMap:  make(map[K]*freqNode[K], capacity),
		listMap:  make(map[int]*list[K], capacity),
		capacity: capacity,
		min:      0,
	}
}

func (s *Store[K]) Inc(key K) {
	node, ok := s.nodeMap[key]
	if !ok {
		s.add(key)
		return
	}

	list, ok := s.listMap[node.freq]
	if ok {
		list.remove(node)
	}

	node.freq++

	nextList, ok := s.listMap[node.freq]
	if !ok {
		nextList = newList[K]()
	}

	nextList.pushBack(node)
	s.listMap[node.freq] = nextList

	if list.size == 0 && s.min == node.freq-1 {
		s.min++
	}
}

func (s *Store[K]) Remove(key K) {
	node, ok := s.nodeMap[key]
	if !ok {
		return
	}

	s.listMap[node.freq].remove(node)
	delete(s.nodeMap, key)
}

func (s *Store[K]) LFU() K {
	minList := s.listMap[s.min]
	leastFrequencyNode := minList.head.next
	key := leastFrequencyNode.key

	return key
}

func (s *Store[K]) Clear() {
	s.nodeMap = make(map[K]*freqNode[K], s.capacity)
	s.listMap = make(map[int]*list[K], s.capacity)
	s.min = 0
}

func (s *Store[K]) add(key K) {
	if _, ok := s.nodeMap[key]; ok {
		s.Inc(key)
		return
	}

	node := &freqNode[K]{key: key, freq: 1}

	s.min = 1
	list, ok := s.listMap[node.freq]
	if !ok {
		list = newList[K]()
	}

	list.pushBack(node)
	s.listMap[node.freq] = list
	s.nodeMap[key] = node
}
