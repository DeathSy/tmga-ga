package main

import (
	"fmt"
	"github.com/deathsy/tmga-ga/services"
	"strconv"
	"time"
)

func init() {
	fmt.Println("Start generating timetable....")
}

func main() {
	ga := services.Genetic{InitialGeneration: 100}
	year := strconv.Itoa(time.Now().Year())
	part := strconv.Itoa(calculatePart())
	ga.Start(part + "/" + year)
}

func calculatePart() int {
	month := time.Now().Month()
	intMonth := 1
	part := 1
	switch {
	case month == time.January:
		intMonth = 1
	case month == time.February:
		intMonth = 2
	case month == time.March:
		intMonth = 3
	case month == time.April:
		intMonth = 4
	case month == time.May:
		intMonth = 5
	case month == time.June:
		intMonth = 6
	case month == time.July:
		intMonth = 7
	case month == time.August:
		intMonth = 8
	case month == time.September:
		intMonth = 9
	case month == time.October:
		intMonth = 10
	case month == time.November:
		intMonth = 11
	case month == time.December:
		intMonth = 12
	}

	switch {
	case intMonth >= 8 && intMonth <= 12:
		part = 2
	case intMonth >= 2 && intMonth <= 5:
		part = 1
	}

	return part
}
