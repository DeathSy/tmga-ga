package models

import "github.com/globalsign/mgo/bson"

type Subject struct {
	Id         bson.ObjectId `bson:"_id" json:"_id"`
	Name       string        `bson:"name" json:"name"`
	Code       string        `bson:"code" json:"code"`
	Student    []string      `bson:"students" json:"students"`
	IsRequired bool          `bson:"isRequired" json:"isRequired"`
}
