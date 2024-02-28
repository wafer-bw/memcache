package volatilelfu

// export for testing.
func (s *Store[K, V]) Capacity() int {
	return s.capacity
}
