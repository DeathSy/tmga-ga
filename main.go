package main

import (
	"github.com/deathsy/tmga-ga/routes"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func init() {
}

func main() {

	r := routes.InitRoutes()

	corsObj := handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(":9000", handlers.CORS(corsObj)(r)))
}
