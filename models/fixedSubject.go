package models

import "github.com/globalsign/mgo/bson"

type FixedSubject struct {
	Id    bson.ObjectId `bson:"_id" json:"_id"`
	Code  string        `bson:"code" json:"code"`
	Name  string        `bson:"name" json:"name"`
	Day   string        `bson:"day" json:"day"`
	Start TimeSlot
	End   TimeSlot
}
