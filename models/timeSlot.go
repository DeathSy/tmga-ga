package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type TimeSlot struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Start     string        `bson:"start" json:"start"`
	End       string        `bson:"end" json:"end"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
