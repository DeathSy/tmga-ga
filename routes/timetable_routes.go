package routes

import (
	"github.com/deathsy/tmga-ga/controllers"
	"github.com/gorilla/mux"
)

func initTimetableRoute(r *mux.Router) {
	s := r.PathPrefix("/timetables").Subrouter()

	s.HandleFunc("", controllers.CreateTimeTable).Methods("POST")
}
