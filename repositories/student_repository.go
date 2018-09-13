package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type StudentRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *StudentRepository) FindAll() ([]models.Student, error) {
	var students []models.Student

	err := r.DB.C(r.Collection).Find(bson.M{}).All(&students)

	return students, err
}

func (r *StudentRepository) Find(id string) (models.Student, error) {
	var student models.Student

	err := r.DB.C(r.Collection).FindId(bson.ObjectIdHex(id)).One(&student)

	return student, err
}
