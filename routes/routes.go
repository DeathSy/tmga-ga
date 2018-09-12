package routes

import (
	"github.com/deathsy/tmga-ga/services"
	"github.com/gorilla/mux"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r2 *http.Request) {
		response := &Response{"Go API server is running", "OK"}
		services.RespondWithJson(w, http.StatusOK, response)
	})

	initTimetableRoute(r)

	return r
}
