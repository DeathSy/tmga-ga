package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type LecturerRepository struct {
	DB *mgo.Database
}

const (
	COLLECTION = "Lecturer"
)

func (m *LecturerRepository) FindAll() ([]models.Lecturer, error) {
	var lecturers []models.Lecturer

	err := m.DB.C(COLLECTION).Find(bson.M{}).All(&lecturers)

	return lecturers, err
}
