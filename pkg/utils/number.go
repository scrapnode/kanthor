package utils

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func AbsInt(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
