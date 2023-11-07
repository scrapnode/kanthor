package validator

import (
	"fmt"
)

func MapNotNil[K comparable, V any](prop string, values map[K]V) Fn {
	return func() error {
		if values == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}
		return nil
	}
}

func MapRequired[K comparable, V any](prop string, values map[K]V) Fn {
	return func() error {
		if err := MapNotNil(prop, values)(); err != nil {
			return err
		}
		if len(values) == 0 {
			return fmt.Errorf("%s contains no item", prop)
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

func SliceRequired[T any](prop string, values []T) Fn {
	return func() error {
		if values == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}

		if len(values) == 0 {
			return fmt.Errorf("%s must not be empty", prop)
		}

		return nil
	}
}
