package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Constraint struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Wants     bool          `bson:"wants" json:"wants"`
	Required  bool          `bson:"required" json:"required"`
	Day       []string      `bson:"day" json:"day"`
	Room      Room          `bson:"room"`
	Subject   Subject       `bson:"subject"`
	StartTime TimeSlot      `bson:"startTime"`
	EndTime   TimeSlot      `bson:"endTime"`
	Lecturer  Lecturer      `bson:"lecturer"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
