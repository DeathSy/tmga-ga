package services

import "github.com/deathsy/tmga-ga/models"

func ExpandingSection(subjectWithSections []models.SubjectSection) []models.Section {
	var sections []models.Section
	for _, subjectSections := range subjectWithSections {
		for _, section := range subjectSections.Sections {
			sections = append(sections, section)
		}
	}

	return sections
}

func ConvertChanToSlice(ch chan Chromosome) []Chromosome {
	r := make([]Chromosome, 0)
	for c := range ch {
		r = append(r, c)
	}

	return r
}
