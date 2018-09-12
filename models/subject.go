package models

import "github.com/globalsign/mgo/bson"

type Subject struct {
	Id                 bson.ObjectId `bson:"id" json:"id"`
	Name               string        `bson:"name" json:"name"`
	Code               string        `bson:"code" json:"code"`
	SectionAmount      int           `bson:"sectionAmount" json:"sectionAmount"`
	StudentsPerSection int           `bson:"studentsPerSection" json:"studentsPerSection"`
	IsCompulsory       bool          `bson:"isCompulsory" json:"isCompulsory"`
}
