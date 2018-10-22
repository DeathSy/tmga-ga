package models

import (
	"github.com/globalsign/mgo/bson"
)

type Constraint struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Wants     bool          `bson:"wants" json:"wants"`
	Room      Room
	Subject   Subject
	StartTime TimeSlot
	EndTime   TimeSlot
	Lecturer  Lecturer
}
