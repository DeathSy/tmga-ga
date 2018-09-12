package models

import "github.com/globalsign/mgo/bson"

type Student struct {
	Id        bson.ObjectId `bson:"id" json:"id"`
	StudentId string        `bson:"studentId" json:"studentId"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt string        `bson:"createdAt" json:"createdAt"`
	UpdatedAt string        `bson:"updatedAt" json:"updatedAt"`
}
