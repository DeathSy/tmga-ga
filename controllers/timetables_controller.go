package controllers

import (
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/repositories"
	"github.com/deathsy/tmga-ga/services"
	"net/http"
)

var timetableRepo repositories.TimetableRepository

func init() {
	session := config.Connect()
	timetableRepo = repositories.TimetableRepository{session, "Timetable"}
}

func GetAllTimetables(w http.ResponseWriter, r *http.Request) {

}

func GetTimetable(w http.ResponseWriter, r *http.Request) {

}

func CreateTimeTable(w http.ResponseWriter, r *http.Request) {
	ga := services.Genetic{InitialGeneration: 100}
	chromosome := ga.Start()

	services.RespondWithJson(w, http.StatusOK, chromosome)
}

func DeleteTimetable(w http.ResponseWriter, r *http.Request) {

}
