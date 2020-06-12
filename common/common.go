package common

func MathMin(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MathMax(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
