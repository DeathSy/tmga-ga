package services

import (
	"fmt"
	"github.com/deathsy/tmga-ga/config"
	"github.com/deathsy/tmga-ga/models"
	"github.com/deathsy/tmga-ga/repositories"
	"math/rand"
	"sort"
	"sync"
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

var fitnessWaitGroup sync.WaitGroup

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	initialReposData()
	timeSlotData, roomData, sectionData, lecturerData = dataPreparation()
	standardGenePattern = sectionPreparation(nil, sectionData)
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

	return geneticFunction(append(roomSlots))
}

func geneticFunction(roomSlots []availableTime) Chromosome {
	var populationPool []Chromosome
	for {
		adam := generateChromosome(roomSlots, sectionPreparation(nil, sectionData))
		eve := generateChromosome(roomSlots, sectionPreparation(nil, sectionData))

		populationPool = sortPopulation(append(populationPool, adam))
		populationPool = sortPopulation(append(populationPool, eve))

		jack, marry := crossover(populationPool[0].Genes, populationPool[1].Genes)

		populationPool = sortPopulation(append(populationPool, jack))
		populationPool = sortPopulation(append(populationPool, marry))

		catherine := mutate(populationPool[0].Genes)

		populationPool = sortPopulation(append(populationPool, catherine))

		fmt.Println("Max fitness", populationPool[0].Fitness)

		if populationPool[0].Fitness > 0.8 {
			break
		}
	}

	fmt.Println("Final fitness", populationPool[0].Fitness)

	return populationPool[0]
}

func generateChromosome(room []availableTime, genes []Gene) Chromosome {
	round := 0
	limit := 0
	for _, gene := range genes {
		limit += gene.Section.Time / 30
	}

	for round < limit {
		roomIndex := rand.Intn(len(room))
		geneIndex := rand.Intn(len(genes))

		if len(genes[geneIndex].TimeSlot)*30 < genes[geneIndex].Section.Time {
			genes[geneIndex].TimeSlot = append(genes[geneIndex].TimeSlot, room[roomIndex])
			genes[geneIndex].TimeSlot = sortTimeSlot(genes[geneIndex].TimeSlot)
			tmpSlice := append(room[:roomIndex])
			copyTmp := make([]availableTime, len(tmpSlice))
			copy(copyTmp, tmpSlice)
			room = append(copyTmp, room[roomIndex+1:]...)
			round += 1
		}
	}
	return calculateFitness(Chromosome{Genes: genes})
}

func calculateFitness(chromosome Chromosome) Chromosome {
	resultCh := make(chan float64, 4)

	fitnessWaitGroup.Add(1)
	go timeCheck(resultCh, transformToTimeBase(chromosome))

	fitnessWaitGroup.Add(1)
	go lecturerCheck(resultCh, transformToLecturerBase(chromosome))

	fitnessWaitGroup.Add(1)
	go consecutiveTimeCheck(resultCh, chromosome.Genes)

	fitnessWaitGroup.Add(1)
	go roomTypeCheck(resultCh, chromosome.Genes)

	fitnessWaitGroup.Wait()
	close(resultCh)

	score := 0.0
	round := 0

	for result := range resultCh {
		score += result
		round++
	}

	chromosome.Fitness = score / float64(round)

	return chromosome
}

func timeCheck(resultCh chan float64, timeBaseChromosome map[string][]models.Section) {
	defer fitnessWaitGroup.Done()

	score := 0.0
	round := 0

	for _, gene := range timeBaseChromosome {
		if len(gene) <= 1 {
			score++
		}
		round++
	}

	resultCh <- score / float64(round)
}

func lecturerCheck(resultCh chan float64, lecturerBaseChromosome map[string][]availableTime) {
	defer fitnessWaitGroup.Done()

	score := 0.0
	round := 0

	for _, gene := range lecturerBaseChromosome {
		for key, geneTime := range gene {
			tmpSlice := append(gene[:key])
			copyTmp := make([]availableTime, len(tmpSlice))
			copy(copyTmp, tmpSlice)
			otherTime := append(copyTmp, gene[key+1:]...)
			check := true
			for _, slot := range otherTime {
				case1 := slot.Room.Name == geneTime.Room.Name
				case2 := slot.Day == geneTime.Day
				case3 := slot.Time.Start == geneTime.Time.Start

				check = check && (case1 || !(case2 && case3))
			}
			if check {
				score++
			}
			round++
		}
	}

	resultCh <- score / float64(round)
}

func consecutiveTimeCheck(resultCh chan float64, standardChromosome []Gene) {
	defer fitnessWaitGroup.Done()

	score := 0.0
	round := 0

	for _, gene := range standardChromosome {
		for key, geneTime := range gene.TimeSlot {
			tmpSlice := append(gene.TimeSlot[:key])
			copyTmp := make([]availableTime, len(tmpSlice))
			copy(copyTmp, tmpSlice)
			otherTime := append(copyTmp, gene.TimeSlot[key+1:]...)
			check := true
			for _, slot := range otherTime {
				case1 := geneTime.Room.Id == slot.Room.Id
				case2 := geneTime.Day == slot.Day
				check = check && (case1 && case2)
			}
			if key+1 < len(gene.TimeSlot) {
				check = geneTime.Time.End == gene.TimeSlot[key+1].Time.Start
			}
			if check {
				score++
			}
			round++
		}
	}

	resultCh <- score / float64(round)
}

func roomTypeCheck(resultCh chan float64, standardChromosome []Gene) {
	defer fitnessWaitGroup.Done()

	score := 0.0
	round := 0

	for _, gene := range standardChromosome {
		for _, geneTime := range gene.TimeSlot {
			if gene.Section.Type == geneTime.Room.Type.Name {
				score++
			}
			round++
		}
	}

	resultCh <- score / float64(round)
}

func crossover(chromosomeA []Gene, chromosomeB []Gene) (Chromosome, Chromosome) {
	randIndex := rand.Intn(len(chromosomeA))

	sliceTmpA := append(chromosomeA[:randIndex])
	copyTmpA := make([]Gene, len(sliceTmpA))
	copy(copyTmpA, sliceTmpA)

	sliceTmpB := append(chromosomeB[:randIndex])
	copyTmpB := make([]Gene, len(sliceTmpB))
	copy(copyTmpB, sliceTmpB)

	jack := append(copyTmpA, chromosomeB[randIndex:]...)
	mary := append(sliceTmpB, chromosomeA[randIndex:]...)

	return calculateFitness(Chromosome{Genes: jack}), calculateFitness(Chromosome{Genes: mary})
}

func mutate(chromosome []Gene) Chromosome {
	mutationPercentage := len(chromosome) * 10 / 100

	var randGeneIndex, randRoomIndex []int
	var limit int

	for mutationPercentage > 0 {
		randNum := rand.Intn(len(chromosome))
		randGeneIndex = append(randGeneIndex, randNum)

		limit += chromosome[randNum].Section.Time / 30

		mutationPercentage--
	}

	for limit > 0 {
		randNum := rand.Intn(len(roomSlots))
		randRoomIndex = append(randRoomIndex, randNum)

		limit--
	}

	tmpSlice := append(chromosome)
	copyTmp := make([]Gene, len(chromosome))
	copy(copyTmp, tmpSlice)

	for _, index := range randGeneIndex {
		copyTmp[index].TimeSlot = []availableTime{}
		for timeLimit := 0; timeLimit < copyTmp[index].Section.Time/30; timeLimit++ {
			copyTmp[index].TimeSlot = append(copyTmp[index].TimeSlot, roomSlots[randRoomIndex[0]])
			randRoomIndex = append(randRoomIndex[1:])
		}
	}

	return calculateFitness(Chromosome{Genes: copyTmp})
}

func transformToTimeBase(chromosome Chromosome) map[string][]models.Section {
	timeMap := convertTimeSliceToMap(map[string][]models.Section{}, timePreparation(nil, roomData, DAYS, timeSlotData))

	for _, gene := range chromosome.Genes {
		for _, slot := range gene.TimeSlot {
			key := slot.Room.Name + slot.Day + slot.Time.Start
			tmpKey := append(timeMap[key])
			timeMap[key] = append(tmpKey, gene.Section)
		}
	}

	return timeMap
}

func transformToLecturerBase(chromosome Chromosome) map[string][]availableTime {
	lecturerMap := convertLecturerSliceToMap(map[string][]availableTime{}, lecturerData)

	for _, gene := range chromosome.Genes {
		for _, lecturer := range gene.Section.Lecturers {
			lecturerMap[lecturer] = append(lecturerMap[lecturer], gene.TimeSlot...)
		}
	}
	return lecturerMap
}

func sortPopulation(population []Chromosome) []Chromosome {
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	return population
}

func sortTimeSlot(timeSlot []availableTime) []availableTime {
	sort.SliceStable(timeSlot, func(i, j int) bool {
		dayI := convertDayToInt(timeSlot[i].Day)
		dayJ := convertDayToInt(timeSlot[j].Day)

		if dayI == dayJ {
			return timeSlot[i].Time.Start < timeSlot[j].Time.Start
		}

		return dayI < dayJ
	})

	return timeSlot
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

func convertDayToInt(day string) int {
	var dayInt int
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
