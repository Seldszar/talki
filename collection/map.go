package collection

import (
	"sync"
)

type Map[K, V comparable] struct {
	mu sync.RWMutex

	data map[K]V
}

func NewMap[K, V comparable]() Map[K, V] {
	return Map[K, V]{
		data: make(map[K]V),
	}
}

func (s *Map[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *Map[K, V]) Get(k K) V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, _ := s.data[k]

	return v
}

func (s *Map[K, V]) Has(k K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[k]

	return ok
}

func (s *Map[K, V]) Set(k K, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[k] = v
}

func (s *Map[K, V]) Delete(k K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, k)
}

func (s *Map[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.data {
		delete(s.data, k)
	}
}

func (s *Map[K, V]) Each(fn func(K, V) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range s.data {
		if fn(k, v) {
			break
		}
	}
}
