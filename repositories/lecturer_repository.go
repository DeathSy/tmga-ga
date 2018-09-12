package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type LecturerRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *LecturerRepository) FindAll() ([]models.Lecturer, error) {
	var lecturers []models.Lecturer

	err := r.DB.C(r.Collection).Find(bson.M{}).All(&lecturers)

	return lecturers, err
}
