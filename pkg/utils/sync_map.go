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
	s.m.Delete(key)
	return
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}
