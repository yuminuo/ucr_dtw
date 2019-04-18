package UcrDtw

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadDataFromFile(t *testing.T) {
	data := ReadDataFromFile("./test_data/Data.txt")
	assert.Equal(t, len(data.items), 1000000)
}

func TestFindSimilar(t *testing.T) {
	len_ := int64(128)
	window := 0.05
	epoch := 100000

	query_data := ReadQueryFromFile("./test_data/Query2.txt")
	data := ReadDataFromFile("./test_data/Data.txt")

	query := PrepareQuery(query_data, len_, window)
	_, loc, _ := FindSimilar(data, query, 3.87, epoch)

	assert.Equal(t, int64(223370), loc)

	query_data = ReadQueryFromFile("./test_data/Query.txt")

	query = PrepareQuery(query_data, len_, window)
	_, loc, _ = FindSimilar(data, query, 3.87, epoch)

	assert.Equal(t, int64(756562), loc)

}
