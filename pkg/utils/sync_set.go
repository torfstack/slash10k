package utils

type SyncSet[T comparable] struct {
	m map[T]struct{}
}

func NewSyncSet[T comparable]() SyncSet[T] {
	return SyncSet[T]{m: make(map[T]struct{})}
}

func (s *SyncSet[T]) Add(value T) {
	s.m[value] = struct{}{}
}

func (s *SyncSet[T]) Remove(value T) {
	delete(s.m, value)
}

func (s *SyncSet[T]) Contains(value T) bool {
	_, ok := s.m[value]
	return ok
}
