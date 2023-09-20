package validator

import "fmt"

// StructNotNil is a predeclared identifier representing the zero value for a pointer, channel, func, interface, map, or slice type.
func MapNotNil[K comparable, V any](prop string, value map[K]V) Fn {
	return func() error {
		if value == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}
		return nil
	}
}

func PointerNotNil[T any](prop string, value *T) Fn {
	return func() error {
		if value == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}
		return nil
	}
}

func SliceRequired[T any](prop string, value []T) Fn {
	return func() error {
		if value == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}

		if len(value) == 0 {
			return fmt.Errorf("%s is required", prop)
		}

		return nil
	}
}
