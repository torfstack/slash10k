package utils

import "sync"

type SyncMap[K, V any] struct {
	m sync.Map
}

func (s *SyncMap[K, V]) LoadAndRemove(key K) (value V, ok bool) {
	v, ok := s.m.Load(key)
	if !ok {
		return
	}

	value, ok = v.(V)
	s.Remove(key)
	return
}

func (s *SyncMap[K, V]) Remove(key K) {
	s.m.Delete(key)
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}

func (s *SyncMap[K, V]) Contains(key K) bool {
	_, ok := s.m.Load(key)
	return ok
}
