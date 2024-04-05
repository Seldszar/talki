package collection

import (
	"sync"
)

type Set[T comparable] struct {
	mu sync.RWMutex

	data map[T]struct{}
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		data: make(map[T]struct{}),
	}
}

func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *Set[T]) Has(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[v]

	return ok
}

func (s *Set[T]) Add(v ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range v {
		s.data[val] = struct{}{}
	}
}

func (s *Set[T]) Delete(v ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range v {
		delete(s.data, val)
	}
}

func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for v := range s.data {
		delete(s.data, v)
	}
}

func (s *Set[T]) Each(fn func(T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for v := range s.data {
		if fn(v) {
			break
		}
	}
}
