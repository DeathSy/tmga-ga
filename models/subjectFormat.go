package models

import "github.com/globalsign/mgo/bson"

type SubjectFormat struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt string        `bson:"createdAt" json:"updatedAt"`
	UpdatedAt string        `bson:"updatedAt" json:"updatedAt"`
}
