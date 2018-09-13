package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type SubjectRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *SubjectRepository) FindAll() ([]models.Subject, error) {
	var subjects []models.Subject

	err := r.DB.C(r.Collection).Find(bson.M{}).All(&subjects)

	return subjects, err
}

func (r *SubjectRepository) Find(id string) (models.Subject, error) {
	var subject models.Subject

	err := r.DB.C(r.Collection).FindId(bson.ObjectIdHex(id)).One(&subject)

	return subject, err
}
