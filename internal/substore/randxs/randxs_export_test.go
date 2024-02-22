package randxs

// export for testing.
func (s *Store[K]) Keys() ([]K, func()) {
	s.mu.Lock()

	return s.keys, s.mu.Unlock
}

// export for testing.
func (s *Store[K]) KeyIndices() (map[K]int, func()) {
	s.mu.Lock()

	return s.keyIndices, s.mu.Unlock
}
