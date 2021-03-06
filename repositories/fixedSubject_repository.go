package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type FixedSubjectRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *FixedSubjectRepository) FindAll() ([]models.FixedSubject, error) {
	var fixedSubjects []models.FixedSubject

	query := []bson.M{
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
				"from":         "TimeSlot",
				"localField":   "startTimeId",
				"foreignField": "_id",
				"as":           "start",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "TimeSlot",
				"localField":   "endTimeId",
				"foreignField": "_id",
				"as":           "end",
			},
		},
		{"$unwind": "$subject"},
		{"$unwind": "$start"},
		{"$unwind": "$end"},
	}

	err := r.DB.C(r.Collection).Pipe(query).All(&fixedSubjects)

	return fixedSubjects, err
}
