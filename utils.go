package UcrDtw

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadQueryFromFile(pth string) []float64 {
	file, err := os.Open(pth)

	if err != nil {
		log.Fatal(err)
	}

	var query []float64
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		for _, v := range strings.Split(row, " ") {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			query = append(query, f)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return query

}

func ReadDataFromFile(pth string) Queue {
	file, err := os.Open(pth)

	if err != nil {
		log.Fatal(err)
	}

	queue := Queue{}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dat, err := ioutil.ReadFile(pth)
	for _, v := range strings.Split(string(dat), " ") {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		queue.Enqueue(f)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return queue

}
