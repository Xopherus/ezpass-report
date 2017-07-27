package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type receipt struct {
	Transponder string
	Cost        float64
}

func sum(x []receipt) float64 {
	var sum float64

	for _, t := range x {
		sum += t.Cost
	}
	return sum
}

func main() {
	fr, err := os.Open("./ezpass.csv")
	if err != nil {
		log.Printf("Could not open file: %s", err)
	}

	r := csv.NewReader(fr)
	data, err := r.ReadAll()
	if err != nil {
		log.Printf("Could not read csv file: %s", err)
	}

	allTolls := make(map[string][]receipt)

	// group tolls by transponder/license #
	for _, line := range data[1:] {
		tr := line[2]

		// skip transactions which don't have transponder/license #
		if tr == "-" {
			continue
		}

		cost, err := strconv.ParseFloat(line[12][2:len(line[12])-1], 64)
		if err != nil {
			log.Printf("Could not convert cost to float: %s", err)
		}

		allTolls[tr] = append(allTolls[tr], receipt{tr, cost})
	}

	totals := make(map[string]float64)

	// sum tolls
	for tr, tolls := range allTolls {
		totals[tr] = sum(tolls)

		log.Printf("Transponder/License [%s]: %.2f", tr, totals[tr])
	}
}
