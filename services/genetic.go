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
	InitialGeneration int
}

type Chromosome struct {
	Genes   []models.Gene
	Fitness float64
}

var DAYS = []string{"MON", "TUE", "WED", "THU", "FRI"}

var sectionData []models.SubjectSection
var timeSlotData []models.TimeSlot
var roomData []models.Room

var timetableRepo repositories.TimetableRepository

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	session := config.Connect()
	timeSlotRepo := repositories.TimeslotRepository{DB: session, Collection: "TimeSlot"}
	roomRepo := repositories.RoomRepository{DB: session, Collection: "Room"}
	sectionRepo := repositories.SubjectSectionRepository{DB: session, Collection: "SubjectSection"}
	timetableRepo = repositories.TimetableRepository{DB: session, Collection: "Timetable"}

	sectionData, _ = sectionRepo.FindAll()
	timeSlotData, _ = timeSlotRepo.FindAll()
	roomData, _ = roomRepo.FindAll()
	sort.SliceStable(timeSlotData, func(i, j int) bool {
		return timeSlotData[i].Start < timeSlotData[j].Start
	})
}

func (g *Genetic) Start(semester string) Chromosome {
	var generateGroup sync.WaitGroup
	populationCh := make(chan Chromosome, g.InitialGeneration)

	for round := 0; round < g.InitialGeneration; round++ {
		generateGroup.Add(1)
		go func() {
			defer generateGroup.Done()

			chromosome := generateChromosome()
			populationCh <- chromosome
		}()
	}
	generateGroup.Wait()
	close(populationCh)

	populationPool := ConvertChanToSlice(populationCh)
	sortPopulation(populationPool)

	timetableRepo.UpdateOrInsert(transformPopulationToTimetable(populationPool[0], semester))

	for populationPool[0].Fitness < 0.85 {
		tmpSlice := append(populationPool[:g.InitialGeneration*20/100])
		copyTmp := make([]Chromosome, len(tmpSlice))
		copy(copyTmp, tmpSlice)
		populationPool = append(copyTmp)

		var tmpPool []Chromosome
		for outerPointer, bestPopulation := range populationPool {
			if outerPointer == len(populationPool)-1 {
				break
			}
			for innerPointer := outerPointer + 1; innerPointer < len(populationPool); innerPointer++ {
				crossoverResult1, crossoverResult2 := crossover(bestPopulation.Genes, populationPool[innerPointer].Genes)
				tmpPool = append(tmpPool, calculateFitness(crossoverResult1))
				tmpPool = append(tmpPool, calculateFitness(crossoverResult2))
			}
		}
		populationPool = append(populationPool, tmpPool...)
		sortPopulation(populationPool)

		tmpSlice = append(populationPool[:g.InitialGeneration-(g.InitialGeneration*20/100)])
		copyTmp = make([]Chromosome, len(tmpSlice))
		copy(copyTmp, tmpSlice)
		populationPool = append(copyTmp)

		first10PercentOfPopulation := tmpSlice[:g.InitialGeneration*10/100]
		var tmpMutationPool []Chromosome
		for round := 0; round < len(first10PercentOfPopulation); round++ {
			mutateChromosome := mutate(first10PercentOfPopulation[round].Genes)
			tmpMutationPool = append(tmpMutationPool, calculateFitness(mutateChromosome))
		}

		populationPool = append(populationPool, tmpMutationPool...)

		var tmpGeneratePool []Chromosome
		for round := 0; round < len(first10PercentOfPopulation); round++ {
			mutateChromosome := generateChromosome()
			tmpGeneratePool = append(tmpGeneratePool, calculateFitness(mutateChromosome.Genes))
		}

		populationPool = append(populationPool, tmpGeneratePool...)
		timetableRepo.UpdateOrInsert(transformPopulationToTimetable(populationPool[0], semester))
		fmt.Println("Fitness", populationPool[0].Fitness)
	}

	return populationPool[0]
}

func transformPopulationToTimetable(chromosome Chromosome, semester string) *models.Timetable {
	return &models.Timetable{FitnessLevel: chromosome.Fitness, Sections: chromosome.Genes, Semester: semester}
}

func generateChromosome() Chromosome {
	var genes []models.Gene
	sections := ExpandingSection(sectionData)

	for _, section := range sections {
		randRoomIndex := rand.Intn(len(roomData))
		randDayIndex := rand.Intn(len(DAYS))
		randTimeSlotIndex := rand.Intn(len(timeSlotData) - section.Time/30)

		genes = append(
			genes,
			models.Gene{
				Section: section,
				Room:    roomData[randRoomIndex],
				Day:     DAYS[randDayIndex],
				Time:    timeSlotData[randTimeSlotIndex : randTimeSlotIndex+section.Time/30]})
	}

	return calculateFitness(genes)
}

func calculateFitness(genes []models.Gene) Chromosome {
	var fitnessGroup sync.WaitGroup
	resultCh := make(chan float64, 3)

	fitnessGroup.Add(1)
	go func() {
		defer fitnessGroup.Done()

		timeBaseChromosome := transformToTimeBaseChromosome(genes)
		score := 0.0
		round := 0

		for _, timeBaseGene := range timeBaseChromosome {
			if len(timeBaseGene) <= 1 {
				score++
			}
			round++
		}

		resultCh <- score / float64(round)
	}()

	fitnessGroup.Add(1)
	go func() {
		defer fitnessGroup.Done()
		score := 0.0
		round := 0

		for _, gene := range genes {
			if gene.Section.Type == gene.Room.Type.Id.Hex() {
				score++
			}
			round++
		}

		resultCh <- score / float64(round)
	}()

	fitnessGroup.Add(1)
	go func() {
		defer fitnessGroup.Done()

		timeBaseFromWithLecturer := transformToTimeBaseWithLecturer(genes)
		score := 0.0
		round := 0

		for _, timeBaseGene := range timeBaseFromWithLecturer {
			if len(timeBaseGene) <= 1 {
				score++
			}
			round++
		}

		resultCh <- score / float64(round)
	}()

	fitnessGroup.Wait()
	close(resultCh)

	finalResult := 0.0
	finalRound := len(resultCh)
	for result := range resultCh {
		finalResult += result
	}

	return Chromosome{Genes: genes, Fitness: finalResult / float64(finalRound)}
}

func crossover(chromosomeA []models.Gene, chromosomeB []models.Gene) ([]models.Gene, []models.Gene) {
	randIndex := rand.Intn(len(chromosomeA))

	aTmp := append(chromosomeA[:randIndex])
	copyATmp := make([]models.Gene, len(aTmp))
	copy(copyATmp, aTmp)
	resultA := append(copyATmp)

	bTmp := append(chromosomeB[:randIndex])
	copyBTmp := make([]models.Gene, len(bTmp))
	copy(copyBTmp, bTmp)
	resultB := append(copyATmp)

	resultA = append(resultA, chromosomeB[randIndex:]...)
	resultB = append(resultB, chromosomeA[randIndex:]...)

	return resultA, resultB
}

func mutate(chromosome []models.Gene) []models.Gene {
	mutateRound := len(chromosome) * 10 / 100
	for index := 0; index < mutateRound; index++ {
		randGeneIndex := rand.Intn(len(chromosome))
		randRoomIndex := rand.Intn(len(roomData))
		randDayIndex := rand.Intn(len(DAYS))
		randTimeSlotIndex := rand.Intn(len(timeSlotData) - chromosome[randGeneIndex].Section.Time/30)

		chromosome[randGeneIndex] = models.Gene{
			Section: chromosome[randGeneIndex].Section,
			Room:    roomData[randRoomIndex],
			Day:     DAYS[randDayIndex],
			Time:    timeSlotData[randTimeSlotIndex : randTimeSlotIndex+chromosome[randGeneIndex].Section.Time/30],
		}
	}
	return chromosome
}

func sortPopulation(populationPool []Chromosome) {
	sort.SliceStable(populationPool, func(i, j int) bool {
		return populationPool[i].Fitness > populationPool[j].Fitness
	})
}

func transformToTimeBaseChromosome(genes []models.Gene) map[string][]models.Section {
	timeMap := make(map[string][]models.Section)
	for _, room := range roomData {
		for _, day := range DAYS {
			for _, slot := range timeSlotData {
				timeMap[room.Name+day+slot.Start] = []models.Section{}
			}
		}
	}

	for _, gene := range genes {
		for _, slot := range gene.Time {
			key := gene.Room.Name + gene.Day + slot.Start
			timeMap[key] = append(timeMap[key], gene.Section)
		}
	}

	return timeMap
}

func transformToTimeBaseWithLecturer(genes []models.Gene) map[string][]string {
	timeMap := make(map[string][]string)
	for _, day := range DAYS {
		for _, slot := range timeSlotData {
			timeMap[day+slot.Start] = []string{}
		}
	}

	for _, gene := range genes {
		for _, timeSlot := range gene.Time {
			for _, lecture := range gene.Section.Lecturers {
				key := gene.Day + timeSlot.Start
				timeMap[key] = append(timeMap[key], lecture)
			}
		}
	}

	return timeMap
}
