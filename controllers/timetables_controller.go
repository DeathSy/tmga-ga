package controllers

import (
	"github.com/deathsy/tmga-ga/services"
	"net/http"
	"os/exec"
)

func CreateTimeTable(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("genetic")
	cmd.Start()

	response := struct {
		Message string
	}{Message: "success"}

	services.RespondWithJson(w, http.StatusOK, response)
}
