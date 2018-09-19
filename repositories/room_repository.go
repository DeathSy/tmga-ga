package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type RoomRepository struct {
	DB         *mgo.Database
	Collection string
}

func (r *RoomRepository) FindAll() ([]models.Room, error) {
	var rooms []models.Room

	err := r.DB.C(r.Collection).Pipe([]bson.M{
		{"$lookup": bson.M{
			"from": "SubjectFormat", "localField": "subjectFormatId", "foreignField": "_id", "as": "format"},
		},
	}).All(&rooms)

	return rooms, err
}

func (r *RoomRepository) Find(id string) (models.Room, error) {
	var room models.Room

	err := r.DB.C(r.Collection).FindId(bson.ObjectIdHex(id)).One(&room)

	return room, err
}
