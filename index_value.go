package UcrDtw

import "math"

type IndexValue struct {
	Value float64
	Index int64
}

type IndexArray []IndexValue

func (arr IndexArray) Len() int { return len(arr) }
func (arr IndexArray) Less(i, j int) bool {
	a := arr[i]
	b := arr[j]
	//val := int16(math.Abs(b.Value) - math.Abs(a.Value))

	return math.Abs(b.Value)-math.Abs(a.Value) > 0
}
func (arr IndexArray) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
