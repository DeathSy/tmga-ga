package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type ConstraintRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *ConstraintRepository) FindAll() ([]models.Constraint, error) {
	var constraints []models.Constraint

	query := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "Room",
				"localField":   "roomId",
				"foreignField": "_id",
				"as":           "room",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "Subject",
				"localField":   "subjectId",
				"foreignField": "_id",
				"as":           "subject",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "Lecturer",
				"localField":   "lecturerId",
				"foreignField": "_id",
				"as":           "lecturer",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "TimeSlot",
				"localField":   "startTimeId",
				"foreignField": "_id",
				"as":           "startTime",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "TimeSlot",
				"localField":   "endTimeId",
				"foreignField": "_id",
				"as":           "endTime",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$room",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$subject",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$lecturer",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$startTime",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$endTime",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"required": false,
			},
		},
	}

	err := r.DB.C(r.Collection).Pipe(query).All(&constraints)

	return constraints, err
}
