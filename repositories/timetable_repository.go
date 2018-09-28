package repositories

import (
	"github.com/deathsy/tmga-ga/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sort"
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

	sort.SliceStable(timetable.Sections, func(i, j int) bool {
		return convertDayToInt(timetable.Sections[i].Day) < convertDayToInt(timetable.Sections[j].Day)
	})

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

func convertDayToInt(day string) int {
	dayInt := 1
	switch {
	case day == "MON":
		dayInt = 1
	case day == "TUE":
		dayInt = 2
	case day == "WED":
		dayInt = 3
	case day == "THU":
		dayInt = 4
	case day == "FRI":
		dayInt = 5
	}

	return dayInt
}
