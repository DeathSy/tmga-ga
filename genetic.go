package main

import (
	"fmt"
	"github.com/deathsy/tmga-ga/services"
)

func main() {
	fmt.Println("Start generating timetable....")

	ga := services.Genetic{InitialGeneration: 100}
	ga.Start()
}
