package UcrDtw

func dist(x float64, y float64) float64 {
	return (x - y) * (x - y)
}

func min_float(x float64, y float64) float64 {
	if x < y {
		return x
	}

	return y
}

func max_int(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
func min_int(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
