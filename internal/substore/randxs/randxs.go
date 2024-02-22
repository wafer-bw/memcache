package randxs

import (
	"math/rand"
	"sync"
)

type Store[K comparable] struct {
	mu         sync.RWMutex
	keys       []K       // permits random key selection
	keyIndices map[K]int // permits fast removal from the keys slice
}

func New[K comparable](startingCapacity int) *Store[K] {
	return &Store[K]{
		keys:       make([]K, 0, startingCapacity),
		keyIndices: make(map[K]int, startingCapacity),
	}
}

func (s *Store[K]) Add(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.keyIndices[key]; ok {
		return
	}

	s.keys = append(s.keys, key)
	s.keyIndices[key] = len(s.keys) - 1
}

func (s *Store[K]) Remove(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.keyIndices[key]
	if !ok {
		return
	}

	delete(s.keyIndices, key)

	isLast := index == len(s.keys)-1
	s.keys[index] = s.keys[len(s.keys)-1]
	s.keys = s.keys[:len(s.keys)-1]
	if !isLast {
		s.keyIndices[s.keys[index]] = index
	}
}

func (s *Store[K]) RandomKey() (K, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.keys) == 0 {
		return *new(K), false
	}

	return s.keys[rand.Intn(len(s.keys))], true
}
