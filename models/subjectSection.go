package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type section struct {
	Name      string   `bson:"name" json:"name"`
	SubjectId string   `bson:"subjectId" json:"subjectId"`
	Type      string   `bson:"type" json:"type"`
	Lecturers []string `bson:"lecturers" json:"lecturers"`
}

type SubjectSection struct {
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	Sections  []section
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
