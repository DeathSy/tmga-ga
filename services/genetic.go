package services

import (
	"fmt"
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/models"
	"github.com/deathsy/tmga-ga/repositories"
	"math/rand"
	"sort"
	"time"
)

type Genetic struct {
}

type Gene struct {
	Section  models.Section
	TimeSlot []availableTime
}

type Chromosome struct {
	Genes   []Gene
	Fitness float64
}

type availableTime struct {
	Day  string          `json:"Day"`
	Room models.Room     `json:"Room"`
	Time models.TimeSlot `json:"timeSlot"`
}

var timeSlotRepo repositories.TimeslotRepository
var lecturerRepo repositories.LecturerRepository
var roomRepo repositories.RoomRepository
var studentRepo repositories.StudentRepository
var sectionRepo repositories.SubjectSectionRepository

var roomData []models.Room
var timeSlotData []models.TimeSlot
var sectionData []models.SubjectSection
var lecturerData []models.Lecturer

var DAYS = []string{"MON", "TUE", "WED", "THU", "FRI"}

var standardGenePattern []Gene
var roomSlots []availableTime

func init() {
	initialReposData()
	timeSlotData, roomData, sectionData, lecturerData = dataPreparation()
}

func initialReposData() {
	session := config.Connect()
	timeSlotRepo = repositories.TimeslotRepository{session, "TimeSlot"}
	lecturerRepo = repositories.LecturerRepository{session, "Lecturer"}
	roomRepo = repositories.RoomRepository{session, "Room"}
	studentRepo = repositories.StudentRepository{session, "Students"}
	sectionRepo = repositories.SubjectSectionRepository{session, "SubjectSection"}
}

func (g *Genetic) Start() Chromosome {
	roomSlots = timePreparation(nil, roomData, DAYS, timeSlotData)
	standardGenePattern = sectionPreparation(nil, sectionData)

	return geneticFunction(nil, roomSlots)
}

func geneticFunction(timetable []Chromosome, roomSlots []availableTime) Chromosome {
	if timetable != nil && timetable[0].Fitness >= 0.98 {
		fmt.Println("Final calculateFitness value", timetable[0].Fitness)
		return timetable[0]
	}

	adam := generateChromosome(Chromosome{}, roomSlots)
	eve := generateChromosome(Chromosome{}, roomSlots)

	timetable = sortPopulation(append(timetable, adam, eve))

	jack, marry := crossover(timetable[0], timetable[1])

	timetable = sortPopulation(append(timetable, jack, marry))

	newJack := mutate(timetable[0], len(standardGenePattern)*10/100)

	timetable = sortPopulation(append(timetable, newJack))

	fmt.Println("Max calculateFitness value", timetable[0].Fitness)

	return geneticFunction(timetable, roomSlots)
}

func generateChromosome(chromosome Chromosome, roomArray []availableTime) Chromosome {
	if len(chromosome.Genes) == 0 {
		chromosome.Genes = append(chromosome.Genes, standardGenePattern...)
		chromosome.Fitness = 0
	}

	sectionIndex := randArrayIndex(len(standardGenePattern))
	roomIndex := randArrayIndex(len(roomArray))
	if len(chromosome.Genes[sectionIndex].TimeSlot)*30 != chromosome.Genes[sectionIndex].Section.Time {
		chromosome.Genes[sectionIndex].TimeSlot = append(chromosome.Genes[sectionIndex].TimeSlot, roomArray[roomIndex])
		roomArray = append(roomArray[:roomIndex], roomArray[roomIndex+1:]...)
	} else {
		for _, gene := range chromosome.Genes {
			if len(gene.TimeSlot)*30 != gene.Section.Time {
				break
			}

			return calculateFitness(chromosome)
		}
	}

	return generateChromosome(chromosome, roomArray)
}

func calculateFitness(chromosome Chromosome) Chromosome {

	standardFunction(chromosome)
	chromosome.Fitness = rand.Float64()

	return chromosome
}

func standardFunction(chromosome Chromosome) Chromosome {
	timeBaseChromosome := transformToTimeBase(
		chromosome,
		convertTimeSliceToMap(map[string][]models.Section{}, roomSlots),
	)

	lecturerBaseChromosome := transformToLectBase(
		chromosome,
		convertLecturerSliceToMap(map[string][]availableTime{}, lecturerData),
	)

	timeScore, timeRound := timeCheck(0.0, 0, timeBaseChromosome)
	lecturerScore, lecturerRound := lecturerCheck(0.0, 0, lecturerBaseChromosome)

	chromosome.Fitness = (timeScore + lecturerScore) / float64(timeRound+lecturerRound)

	return chromosome
}

func timeCheck(score float64, round int, timeBaseChromosome map[string][]models.Section) (float64, int) {
	for _, timeSlot := range timeBaseChromosome {
		round += 1
		if len(timeSlot) <= 1 {
			score += 1
		}
	}

	return score, round
}

func lecturerCheck(score float64, round int, lecturerBaseChromosome map[string][]availableTime) (float64, int) {
	for _, gene := range lecturerBaseChromosome {
		for index, time := range gene {
			otherSlot := append(gene[:index], gene[index+1:]...)
			check := true
			for _, slot := range otherSlot {
				case1 := time.Time.Start == slot.Time.Start
				case2 := time.Day == slot.Day
				case3 := time.Room.Id == slot.Room.Id

				check = check && !(case1 && case2 && case3)
			}

			if check {
				score += 1
			}
			round += 1
		}
	}

	return score, round
}

func transformToTimeBase(chromosome Chromosome, timeMap map[string][]models.Section) map[string][]models.Section {
	for _, gene := range chromosome.Genes {
		for _, time := range gene.TimeSlot {
			key := time.Room.Name + time.Day + time.Time.Start
			timeMap[key] = append(timeMap[key], gene.Section)
		}
	}

	return timeMap
}

func transformToLectBase(chromosome Chromosome, lecturerMap map[string][]availableTime) map[string][]availableTime {
	for _, gene := range chromosome.Genes {
		for _, lecturer := range gene.Section.Lecturers {
			lecturerMap[lecturer] = append(lecturerMap[lecturer], gene.TimeSlot...)
		}
	}

	return lecturerMap
}

func crossover(adam Chromosome, eve Chromosome) (Chromosome, Chromosome) {
	rand.Seed(time.Now().UTC().UnixNano())
	crossingIndex := randArrayIndex(len(standardGenePattern))

	jack := Chromosome{nil, 0}
	marry := Chromosome{nil, 0}

	jack.Genes = append(adam.Genes[:crossingIndex], eve.Genes[crossingIndex:]...)
	marry.Genes = append(eve.Genes[:crossingIndex], adam.Genes[crossingIndex:]...)

	return calculateFitness(jack), calculateFitness(marry)
}

func mutate(chromosome Chromosome, round int) Chromosome {
	if round == 0 {
		return calculateFitness(chromosome)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	mutationIndex := randArrayIndex(len(standardGenePattern))
	chromosome.Genes[mutationIndex].TimeSlot = []availableTime{}
	chromosome.Genes[mutationIndex] = renewGene(chromosome.Genes[mutationIndex])

	return mutate(chromosome, round-1)
}

func renewGene(gene Gene) Gene {
	if len(gene.TimeSlot)*30 == gene.Section.Time {
		return gene
	}

	roomIndex := randArrayIndex(len(roomSlots))
	gene.TimeSlot = append(gene.TimeSlot, roomSlots[roomIndex])

	return renewGene(gene)
}

func randArrayIndex(arraySize int) int {
	return rand.Intn(arraySize)
}

func sortPopulation(population []Chromosome) []Chromosome {
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	return population
}

func convertTimeSliceToMap(result map[string][]models.Section, time []availableTime) map[string][]models.Section {
	if len(time) == 0 {
		return result
	}

	mapKey := time[0].Room.Name + time[0].Day + time[0].Time.Start
	result[mapKey] = []models.Section{}

	return convertTimeSliceToMap(result, time[1:])
}

func convertLecturerSliceToMap(result map[string][]availableTime, lecturers []models.Lecturer) map[string][]availableTime {
	if len(lecturers) == 0 {
		return result
	}

	mapKey := lecturers[0].Id.String()
	result[mapKey] = []availableTime{}

	return convertLecturerSliceToMap(result, lecturers[1:])
}

func dataPreparation() ([]models.TimeSlot, []models.Room, []models.SubjectSection, []models.Lecturer) {
	t, _ := timeSlotRepo.FindAll()
	r, _ := roomRepo.FindAll()
	sec, _ := sectionRepo.FindAll()
	l, _ := lecturerRepo.FindAll()

	return t, r, sec, l
}

func sectionPreparation(genes []Gene, allSections []models.SubjectSection) []Gene {
	genes = reorderSec(genes, allSections[0].Sections)

	if len(allSections) == 1 {
		return genes
	}

	return sectionPreparation(genes, allSections[1:])
}

func reorderSec(genes []Gene, sections []models.Section) []Gene {
	genes = append(genes, Gene{
		sections[0],
		[]availableTime{},
	})

	if len(sections) == 1 {
		return genes
	}

	return reorderSec(genes, sections[1:])
}

func timePreparation(time []availableTime, rooms []models.Room, days []string, slots []models.TimeSlot) []availableTime {
	time = bindRoomWIthDay(time, rooms[0], days, slots)

	if len(rooms) == 1 {
		return time
	}

	return timePreparation(time, rooms[1:], days, slots)
}

func bindRoomWIthDay(time []availableTime, rooms models.Room, days []string, slots []models.TimeSlot) []availableTime {

	time = bindRoomAndDayWithTime(time, rooms, days[0], slots)

	if len(days) == 1 {
		return time
	}

	return bindRoomWIthDay(time, rooms, days[1:], slots)

}

func bindRoomAndDayWithTime(time []availableTime, room models.Room, day string, slots []models.TimeSlot) []availableTime {

	time = append(time, availableTime{day, room, slots[0]})

	if len(slots) == 1 {
		return time
	}

	return bindRoomAndDayWithTime(time, room, day, slots[1:])
}
