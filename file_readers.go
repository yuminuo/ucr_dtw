package UcrDtw

import (
	"bufio"
	"fmt"
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
	splited := strings.Split(string(dat), " ")
	for i, v := range splited {
		if i == len(splited)-1 {
			v = strings.Replace(v, "\r\n", "", -1)
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			if i == len(splited)-1 {
				fmt.Print(v)
			}
			//
			continue
		}
		queue.Enqueue(f)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return queue

}
