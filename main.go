package main

import (
	"github.com/deathsy/tmga-ga/routes"
	"log"
	"net/http"
)

func init() {
}

func main() {

	r := routes.InitRoutes()

	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatal(err)
	}
}
