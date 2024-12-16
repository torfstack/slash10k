package utils

type SyncSet[T comparable] struct {
	m SyncMap[T, struct{}]
}

func NewSyncSet[T comparable]() SyncSet[T] {
	return SyncSet[T]{m: SyncMap[T, struct{}]{}}
}

func (s *SyncSet[T]) Add(value T) {
	s.m.Store(value, struct{}{})
}

func (s *SyncSet[T]) Remove(value T) {
	s.m.Remove(value)
}

func (s *SyncSet[T]) Contains(value T) bool {
	return s.m.Contains(value)
}
