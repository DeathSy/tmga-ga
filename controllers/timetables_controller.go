package controllers

import (
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/repositories"
	"github.com/deathsy/tmga-ga/services"
	"github.com/gorilla/mux"
	"net/http"
	"os/exec"
)

func AllTimetable(w http.ResponseWriter, r *http.Request) {
	session := config.Connect()
	repository := repositories.TimetableRepository{DB: session, Collection: "Timetable"}

	timetables, err := repository.FindAll()
	if err != nil {
		services.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		services.RespondWithJson(w, http.StatusOK, timetables)
	}
}

func FindTimetable(w http.ResponseWriter, r *http.Request) {
	session := config.Connect()
	repository := repositories.TimetableRepository{DB: session, Collection: "Timetable"}

	vars := mux.Vars(r)
	query := r.URL.Query()
	timetable, err := repository.Find(vars["part"] + "/" + vars["year"])
	if err != nil {
		services.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		if query.Get("fitnessLevel") != "true" {
			services.RespondWithJson(w, http.StatusOK, timetable)
		} else {
			services.RespondWithJson(w, http.StatusOK, struct {
				FitnessLevel float64 `json:"fitnessLevel"`
			}{FitnessLevel: timetable.FitnessLevel})
		}
	}
}

func CreateTimeTable(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("genetic")
	cmd.Start()

	response := struct {
		Message string
	}{Message: "success"}

	services.RespondWithJson(w, http.StatusOK, response)
}
