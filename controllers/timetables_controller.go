package controllers

import (
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/repositories"
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

}

func DeleteTimetable(w http.ResponseWriter, r *http.Request) {

}
