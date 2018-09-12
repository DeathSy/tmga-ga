package models

import "github.com/globalsign/mgo/bson"

type Timetable struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`
}
