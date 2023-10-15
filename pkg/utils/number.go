package utils

func MaxInt[T int | int32 | int64](x, y T) T {
	if x < y {
		return y
	}
	return x
}

func MinInt[T int | int32 | int64](x, y T) T {
	if x > y {
		return y
	}
	return x
}

func AbsInt[T int | int32 | int64](x T) T {
	if x > 0 {
		return x
	}
	return -x
}
