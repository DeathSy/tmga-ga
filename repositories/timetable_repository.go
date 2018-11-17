package repositories

import (
	"fmt"
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

	for sectionIndex, section := range timetable.Sections {
		var subject models.Subject
		e := r.DB.C("Subject").FindId(bson.ObjectIdHex(section.Section.SubjectId)).One(&subject)

		for _, lecturer := range section.Section.Lecturers {
			var lect models.Lecturer

			e := r.DB.C("Lecturer").FindId(bson.ObjectIdHex(lecturer)).One(&lect)

			if e != nil {
				panic(e)
			}
			timetable.Sections[sectionIndex].Section.LecturerData = append(timetable.Sections[sectionIndex].Section.LecturerData, lect)
		}

		if e != nil {
			panic(e.Error())
		}
		timetable.Sections[sectionIndex].Section.Subject = subject
	}
	return timetable, err
}

func (r *TimetableRepository) UpdateOrInsert(timetable *models.Timetable) {
	fmt.Println("Fitness", timetable.FitnessLevel)

	query := bson.M{"semester": timetable.Semester}
	_, err := r.DB.C(r.Collection).Upsert(
		query,
		bson.M{
			"$set": bson.M{
				"fitnessLevel": timetable.FitnessLevel,
				"Sections":     timetable.Sections,
			},
		},
	)

	if err != nil {
		panic(err.Error())
	}
}
