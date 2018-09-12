package models

import "github.com/globalsign/mgo/bson"

type Room struct {
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	Capacity  int           `bson:"capacity" json:"capacity"`
	CreatedAt string        `bson:"createdAt" json:"createdAt"`
	UpdatedAt string        `bson:"updatedAt" json:"updatedAt"`
	Type      *SubjectFormat
}
