package safe

import "sync"

type Map[T any] struct {
	mu     sync.Mutex
	data   map[string]T
	sample T
}

func (sm *Map[T]) int() {
	if sm.data == nil {
		sm.data = map[string]T{}
	}
}

func (sm *Map[T]) Set(key string, value T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.int()

	sm.data[key] = value
	sm.sample = value
}

func (sm *Map[T]) Get(key string) (T, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.int()

	value, ok := sm.data[key]
	return value, ok
}

func (sm *Map[T]) Merge(values map[string]T) {
	for k, v := range values {
		sm.Set(k, v)
	}
}

func (sm *Map[T]) Sample() T {
	return sm.sample
}

func (sm *Map[T]) Count() int {
	return len(sm.data)
}

func (sm *Map[T]) Data() map[string]T {
	return sm.data
}

func (sm *Map[T]) Keys() []string {
	keys := []string{}
	if sm.data != nil {
		for key := range sm.data {
			keys = append(keys, key)
		}
	}
	return keys
}
