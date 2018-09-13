package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TimeslotRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *TimeslotRepository) FindAll() ([]models.TimeSlot, error) {
	var timeslots []models.TimeSlot

	err := r.DB.C(r.Collection).Find(bson.M{}).All(&timeslots)

	return timeslots, err
}

func (r *TimeslotRepository) Find(id string) (models.TimeSlot, error) {
	var timeslot models.TimeSlot

	err := r.DB.C(r.Collection).FindId(bson.ObjectIdHex(id)).One(&timeslot)

	return timeslot, err
}
