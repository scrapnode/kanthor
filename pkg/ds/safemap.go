package ds

import "sync"

type SafeMap[T any] struct {
	mu   sync.Mutex
	data map[string]T
}

func (sm *SafeMap[T]) int() {
	if sm.data == nil {
		sm.data = map[string]T{}
	}
}

func (sm *SafeMap[T]) Set(key string, value T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.int()

	sm.data[key] = value
}

func (sm *SafeMap[T]) Get(key string) (T, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.int()

	value, ok := sm.data[key]
	return value, ok
}

func (sm *SafeMap[T]) Count() int {
	return len(sm.data)
}

func (sm *SafeMap[T]) Data() map[string]T {
	return sm.data
}
