package utils

func SliceMerge[T any](dest map[string]T, sources ...map[string]T) {
	for _, source := range sources {
		for key, value := range source {
			dest[key] = value
		}
	}
}
