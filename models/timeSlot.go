package models

import "github.com/globalsign/mgo/bson"

type TimeSlot struct {
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	Start     string        `bson:"start" json:"start"`
	End       string        `bson:"end" json:"end"`
	CreatedAt string        `bson:"createdAt" json:"createdAt"`
	UpdatedAt string        `bson"updatedAt" json:"updatedAt"`
}
