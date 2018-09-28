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

func (r *TimetableRepository) Find(semester string) (models.Timetable, error) {
	var timetable models.Timetable

	err := r.DB.C(r.Collection).Find(bson.M{"semester": semester}).One(&timetable)

	return timetable, err
}

func (r *TimetableRepository) Create(timetable *models.Timetable) (bson.ObjectId, error) {
	objectId := bson.NewObjectId()
	timetable.Id = objectId
	err := r.DB.C(r.Collection).Insert(&timetable)

	return objectId, err
}

func (r *TimetableRepository) Update(timetable *models.Timetable) error {
	query := bson.M{"semester": timetable.Semester}
	err := r.DB.C(r.Collection).Update(query, &timetable)

	return err
}
