package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type FixedSubjectRepository struct {
	DB *mgo.Database
	Collection string
}

func (r *FixedSubjectRepository) findAll() ([]models.FixedSubject, error) {
	var fixedSubjects []models.FixedSubject

	err := r.DB.C(r.Collection).Find(bson.M{}).All(fixedSubjects)

	return fixedSubjects, err
}
