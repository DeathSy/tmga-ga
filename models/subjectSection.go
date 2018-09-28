package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Section struct {
	Name         string   `bson:"name" json:"name"`
	SubjectId    string   `bson:"subjectId" json:"subjectId"`
	Type         string   `bson:"type" json:"type"`
	Lecturers    []string `bson:"lecturers" json:"lecturers"`
	LecturerData []Lecturer
	Subject      Subject
	Time         int `bson:"time" json:"time"`
}

type SubjectSection struct {
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	Sections  []Section
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
