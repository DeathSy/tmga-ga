package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type SubjectSectionRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *SubjectSectionRepository) FindAll() ([]models.SubjectSection, error) {
	var subjectSection []models.SubjectSection

	query := []bson.M{{
		"$lookup": bson.M{
			"from": "Subject",
			"localField": "subjectId",
			"foreignField": "_id",
			"as": "subject",
		},
	},}

	err := r.DB.C(r.Collection).Pipe(query).All(&subjectSection)

	return subjectSection, err
}

func (r *SubjectSectionRepository) Find(id string) (models.SubjectSection, error) {
	var subjectSection models.SubjectSection

	err := r.DB.C(r.Collection).FindId(bson.ObjectIdHex(id)).One(&subjectSection)

	return subjectSection, err
}
