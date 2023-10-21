package memcache

import (
	"container/list"
)

type lruEvictor[K comparable] struct {
	capacity int
	list     *list.List
	elements map[K]*list.Element
}

func (p lruEvictor[K]) Over() int {
	size := p.list.Len()
	if size > p.capacity {
		return size - p.capacity
	}
	return 0
}

func (p lruEvictor[K]) Add(key K) {
	e := p.list.PushFront(key)
	p.elements[key] = e
}

func (p lruEvictor[K]) Use(key K) {
	e := p.elements[key]
	p.list.MoveToFront(e)
}

func (p lruEvictor[K]) Remove(key K) {
	e, ok := p.elements[key]
	if !ok {
		return
	}

	p.list.Remove(e)
	delete(p.elements, key)
}

func (p lruEvictor[K]) Evict() (K, bool) {
	e := p.list.Back()
	if e == nil {
		var k K
		return k, false
	}

	v := p.list.Remove(e)
	delete(p.elements, v.(K))

	return v.(K), true
}
