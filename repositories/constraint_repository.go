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

	err := r.DB.C(r.Collection).Find(bson.M{}).All(constraints)

	return constraints, err
}
