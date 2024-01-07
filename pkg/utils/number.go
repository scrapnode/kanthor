package utils

func Max[T int | int32 | int64](x, y T) T {
	if x < y {
		return y
	}
	return x
}

func Min[T int | int32 | int64](x, y T) T {
	if x > y {
		return y
	}
	return x
}
