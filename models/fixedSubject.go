package models

import "github.com/globalsign/mgo/bson"

type FixedSubject struct {
	Id      bson.ObjectId `bson:"_id" json:"_id"`
	Subject Subject
	Day     string `bson:"day" json:"day"`
	Start   TimeSlot
	End     TimeSlot
}
