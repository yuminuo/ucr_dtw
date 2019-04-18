package UcrDtw

import (
	"fmt"
	"math"
	"testing"
)

func TestDTW_similiarity(t *testing.T) {
	query_data := ReadQueryFromFile("./test_data/Query.txt")
	data := ReadDataFromFile("./test_data/Data.txt")
	var len int64 = 128
	window := 0.05
	epoch := 100000
	query := prepareQuery(query_data, len, window, epoch)
	_, loc, val := similiarity_finder(data, query, 4, epoch)

	fmt.Printf("Location: %v\n", loc)
	fmt.Printf("Value: %v\n", math.Sqrt(val))

}
