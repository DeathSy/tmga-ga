package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type LecturerRepository struct {
	Server   string
	Database string
	Username string
	Password string
}

var db *mgo.Database

const (
	COLLECTION = "lecturer"
)

func (m *LecturerRepository) FindAll() ([]models.Lecturer, error) {
	var lecturers []models.Lecturer

	err := db.C(COLLECTION).Find(bson.M{}).All(&lecturers)

	return lecturers, err
}
