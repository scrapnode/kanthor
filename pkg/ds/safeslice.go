package ds

import "sync"

type SafeSlice[T any] struct {
	mu   sync.Mutex
	data []T
}

func (sm *SafeSlice[T]) Append(values ...T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.data = append(sm.data, values...)
}

func (sm *SafeSlice[T]) Count() int {
	return len(sm.data)
}

func (sm *SafeSlice[T]) Data() []T {
	return sm.data
}
