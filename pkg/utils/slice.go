package utils

func SliceMerge[T any](dest map[string]T, sources ...map[string]T) {
	for _, source := range sources {
		for key, value := range source {
			dest[key] = value
		}
	}
}

func ChunkNext[T int | int32 | int64](prev, end, step T) T {
	if prev+step > end {
		return end
	}
	return prev + step
}
