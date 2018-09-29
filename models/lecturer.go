package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Lecturer struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
