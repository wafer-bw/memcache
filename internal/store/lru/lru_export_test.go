package lru

import "container/list"

// export for testing.
func (s *Store[K, V]) Elements() (map[K]*list.Element, func()) {
	s.mu.Lock()

	return s.elements, s.mu.Unlock
}

// export for testing.
func (s *Store[K, V]) List() (*list.List, func()) {
	s.mu.Lock()

	return s.list, s.mu.Unlock
}
