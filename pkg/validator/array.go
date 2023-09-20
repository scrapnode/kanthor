package validator

func Array[T any](items []T, fn func(i int, item T) error) Fn {
	return func() error {
		for i, item := range items {
			if err := fn(i, item); err != nil {
				return err
			}
		}

		return nil
	}
}
