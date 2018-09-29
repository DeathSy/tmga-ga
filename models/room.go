package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Room struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	Capacity  int           `bson:"capacity" json:"capacity"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	Type      SubjectFormat
}
