package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TimetableRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *TimetableRepository) FindAll() ([]models.Timetable, error) {
	var timetables []models.Timetable

	err := r.DB.C(r.Collection).Find(bson.M{}).All(&timetables)

	return timetables, err

}
