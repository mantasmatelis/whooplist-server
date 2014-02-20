package main

import (
	//"../../whooplist"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type item struct {
	rank       int
	list_id    int
	factual_id string
}

func main() {

	fin := os.Stdin
	bufr := bufio.NewReader(fin)
	reader := csv.NewReader(bufr)

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	/* Chop off first column and first row */
	data = data[1:]
	for i := 0; i < len(data); i++ {
		data[i] = data[i][1:]
	}

	/* Separate header from data */
	header := data[0]
	data = data[1:]

	items := make([]item, 0, 100)

	for _, row := range data {
		rank, err := strconv.Atoi(row[0])
		if err != nil {
			//fmt.Printf("invalid row: %s\n", row[0])
			continue
		}
		for i := 0; i < len(row); i++ {
			if i == 0 {
				continue
			}
			header_val, err := strconv.Atoi(header[i])
			if err != nil {
				//fmt.Printf("invalid header: %s\n", header[i])
				continue
			}
			if row[i] == "" {
				continue
			}
			items = append(items, item{rank: rank, list_id: header_val, factual_id: row[i]})
		}
	}

	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}
