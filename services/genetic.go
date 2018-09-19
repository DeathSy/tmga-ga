package services

import (
	"fmt"
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/models"
	"github.com/deathsy/tmga-ga/repositories"
	"math/rand"
	"sort"
)

type Genetic struct {
}

type Gene struct {
	Section availableSec
}

type Chromosome struct {
	Genes   []Gene
	Fitness float64
}

type availableTime struct {
	day  string
	room models.Room
	time models.TimeSlot
}

type availableSec struct {
	Name      string
	SubjectId string
	Type      string
	Lecturers []string
	TimeUsed  int
	TimeSlot  []availableTime
}

var timeSlotRepo repositories.TimeslotRepository
var lecturerRepo repositories.LecturerRepository
var roomRepo repositories.RoomRepository
var studentRepo repositories.StudentRepository
var sectionRepo repositories.SubjectSectionRepository

var roomData []models.Room
var timeSlotData []models.TimeSlot
var sectionData []models.SubjectSection

var roomSlots []availableTime
var sections []availableSec

var DAYS = []string{"MON", "TUE", "WED", "THU", "FRI"}

func init() {
	initialReposData()
	timeSlotData, roomData, sectionData = dataPreparation()
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
	sections = sectionPreparation(nil, sectionData)

	adam := fitness(generateChromosome(nil, roomSlots, sections))
	eve := fitness(generateChromosome(nil, roomSlots, sections))

	fmt.Println("adam gene:", adam.Genes[0])
	fmt.Println("=======================")
	fmt.Println("eve gene:", eve.Genes[0])
	fmt.Println("=======================")

	perfectChromosome := geneticFunction(adam, eve)

	return perfectChromosome
}

func geneticFunction(chromosome1 Chromosome, chromosome2 Chromosome) Chromosome {
	chromosomeIndex := randArrayIndex(len(chromosome1.Genes))

	crossoveredChro1 := fitness(append(chromosome2.Genes[:chromosomeIndex], chromosome1.Genes[chromosomeIndex+1:]...))
	crossoveredChro2 := fitness(append(chromosome1.Genes[:chromosomeIndex], chromosome2.Genes[chromosomeIndex+1:]...))

	population := []Chromosome{chromosome1, chromosome2, crossoveredChro1, crossoveredChro2}

	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	if population[0].Fitness > 0.8 {
		return population[0]
	}

	round := len(population[0].Genes) * 10 / 100
	mutatedChromosome := fitness(mutate(population[0], round).Genes)

	population = append(population, mutatedChromosome)

	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	if population[0].Fitness > 0.8 {
		return population[0]
	}

	return geneticFunction(population[0], population[1])
}

func fitness(genes []Gene) Chromosome {
	chromosome := Chromosome{genes, 0}

	// todo implement fitness function

	return chromosome
}

func mutate(chromosome Chromosome, round int) Chromosome {
	// todo implement mutation function

	if round == 0 {
		return chromosome
	}

	return mutate(chromosome, round-1)
}

func dataPreparation() ([]models.TimeSlot, []models.Room, []models.SubjectSection) {
	t, _ := timeSlotRepo.FindAll()
	r, _ := roomRepo.FindAll()
	sec, _ := sectionRepo.FindAll()

	return t, r, sec
}

func generateChromosome(gene []Gene, roomArray []availableTime, sectionArray []availableSec) []Gene {

	if len(gene) == len(sections) {
		return gene
	}

	var roomIndex int = randArrayIndex(len(roomArray))
	var sectionIndex int = randArrayIndex(len(sectionArray))

	if len(sectionArray[sectionIndex].TimeSlot)*30 == sectionArray[sectionIndex].TimeUsed {
		gene = append(gene, Gene{sections[sectionIndex]})
		return generateChromosome(
			gene,
			append(roomArray[:roomIndex], roomArray[roomIndex+1:]...),
			append(sectionArray[:sectionIndex], sectionArray[sectionIndex+1:]...),
		)
	}

	sectionArray[sectionIndex].TimeSlot = append(sectionArray[sectionIndex].TimeSlot, roomSlots[roomIndex])

	return generateChromosome(gene, append(roomArray[:roomIndex], roomArray[roomIndex+1:]...), sectionArray)

}

func randArrayIndex(arraySize int) int {
	return rand.Intn(arraySize)
}

func sectionPreparation(secs []availableSec, allSections []models.SubjectSection) []availableSec {
	secs = reorderSec(secs, allSections[0].Sections)

	if len(allSections) == 1 {
		return secs
	}

	return sectionPreparation(secs, allSections[1:])
}

func reorderSec(secs []availableSec, sections []models.Section) []availableSec {
	secs = append(secs, availableSec{
		sections[0].Name,
		sections[0].SubjectId,
		sections[0].Type,
		sections[0].Lecturers,
		sections[0].Time,
		[]availableTime{},
	})

	if len(sections) == 1 {
		return secs
	}

	return reorderSec(secs, sections[1:])
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
