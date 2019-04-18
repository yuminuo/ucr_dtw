package UcrDtw

import (
	"math"
	"math/rand"
)

func QuickSort(slice []IndexValue) []IndexValue {
	length := len(slice)

	if length <= 1 {
		sliceCopy := make([]IndexValue, length)
		copy(sliceCopy, slice)
		return sliceCopy
	}

	m := slice[rand.Intn(length)]

	less := make([]IndexValue, 0, length)
	middle := make([]IndexValue, 0, length)
	more := make([]IndexValue, 0, length)

	for _, item := range slice {
		v := Compat(&item, &m)
		switch {
		case v < 0:
			less = append(less, item)
		case v == 0:
			middle = append(middle, item)
		case v > 0:
			more = append(more, item)
		}
	}

	less, more = QuickSort(less), QuickSort(more)

	less = append(less, middle...)
	less = append(less, more...)

	return less
}

func Compat(a *IndexValue, b *IndexValue) float64 {
	val := math.Abs(b.Value) - math.Abs(a.Value)
	return val
}
