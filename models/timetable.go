package models

import (
	"github.com/globalsign/mgo/bson"
)

type Gene struct {
	Section Section
	Room    Room
	Day     string `bson:"day" json:"day"`
	Time    []TimeSlot
}

type Timetable struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Semester     string        `bson:"semester" json:"semester"`
	Sections     []Gene        `bson:"Sections" json:"Sections"`
	FitnessLevel float64       `bson:"fitnessLevel" json:"fitnessLevel"`
}
