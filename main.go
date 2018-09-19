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

	log.Fatal(http.ListenAndServe(":9000", r))
}
