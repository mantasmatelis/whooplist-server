package main

import (
	"../../whooplist"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type item struct {
	rank       int
	list_id    int
	factual_id string
}

func e(str string) string {
	return strings.Replace(str, "'", "\\'", -1)
}

func main() {
	whooplist.Initialize()
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
			continue
		}
		for i := 0; i < len(row); i++ {
			if i == 0 {
				continue
			}
			header_val, err := strconv.Atoi(header[i])
			if err != nil {
				continue
			}
			if row[i] == "" {
				continue
			}
			items = append(items, item{rank: rank, list_id: header_val, factual_id: row[i]})
		}
	}

	for _, item := range items {
		place, err := whooplist.GetPlaceFactual(item.factual_id)
		if err != nil {
			fmt.Errorf("error: " + err.Error())
		}
		if place == nil {
			place, err = whooplist.FactualPlace(item.factual_id)
			if err != nil {
				fmt.Errorf("error: " + err.Error())
				continue
			}
		}

		fmt.Printf("INSERT INTO wl.place (latitude, longitude, factual_id, name, "+
			"address, locality, region, postcode, country, telephone, "+
			"website, email) "+
			"VALUES (%v, %v, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING id;\n",
			place.Latitude, place.Longitude, e(place.FactualId),
			e(place.Name), e(place.Address), e(place.Locality),
			e(place.Region), e(place.Postcode), e(place.Country),
			e(place.Tel), e(place.Website), e(place.Email))
	}
}
