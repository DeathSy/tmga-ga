package services

import (
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/repositories"
)

var timeSlotRepo repositories.TimeslotRepository
var lecturerRepo repositories.LecturerRepository
var roomRepo repositories.RoomRepository
var studentRepo repositories.StudentRepository

func init() {
	session := config.Connect()
	timeSlotRepo = repositories.TimeslotRepository{session, "timeslots"}
	lecturerRepo = repositories.LecturerRepository{session, "lecturers"}
	roomRepo = repositories.RoomRepository{session, "rooms"}
	studentRepo = repositories.StudentRepository{session, "students"}
}

func main() {

}
