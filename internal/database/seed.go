package database

import (
	"farm/internal/models"
	"log"
)

func s(v string) *string { return &v }

func Seed() {
	seedSpecies()
}

func seedSpecies() {
	var count int64
	DB.Model(&models.Species{}).Count(&count)
	if count > 0 {
		return }

	species := []models.Species{
		{Name: "Cattle"},
		{Name: "Goat"},
		{Name: "Chicken"},
	}
	DB.Create(&species)
	log.Println("Seeded species")

	seedBreeds(species)
	seedVaccines(species)
}

func seedBreeds(species []models.Species) {
	speciesMap := map[string]uint{}
	for _, s := range species {
		speciesMap[s.Name] = s.ID
	}

	breeds := []models.Breed{
		{SpeciesID: speciesMap["Cattle"], Name: "Holstein"},
		{SpeciesID: speciesMap["Cattle"], Name: "Jersey"},
		{SpeciesID: speciesMap["Cattle"], Name: "Angus"},
		{SpeciesID: speciesMap["Cattle"], Name: "Hereford"},
		{SpeciesID: speciesMap["Cattle"], Name: "Brahman"},
		{SpeciesID: speciesMap["Cattle"], Name: "Charolais"},
		{SpeciesID: speciesMap["Cattle"], Name: "Simmental"},
		{SpeciesID: speciesMap["Cattle"], Name: "Limousin"},

		{SpeciesID: speciesMap["Goat"], Name: "Boer"},
		{SpeciesID: speciesMap["Goat"], Name: "Saanen"},
		{SpeciesID: speciesMap["Goat"], Name: "Alpine"},
		{SpeciesID: speciesMap["Goat"], Name: "Nubian"},
		{SpeciesID: speciesMap["Goat"], Name: "LaMancha"},
		{SpeciesID: speciesMap["Goat"], Name: "Toggenburg"},

		{SpeciesID: speciesMap["Chicken"], Name: "Rhode Island Red"},
		{SpeciesID: speciesMap["Chicken"], Name: "Leghorn"},
		{SpeciesID: speciesMap["Chicken"], Name: "Plymouth Rock"},
		{SpeciesID: speciesMap["Chicken"], Name: "Wyandotte"},
		{SpeciesID: speciesMap["Chicken"], Name: "Orpington"},
		{SpeciesID: speciesMap["Chicken"], Name: "Sussex"},
	}
	DB.Create(&breeds)
	log.Printf("Seeded %d breeds", len(breeds))
}

func seedVaccines(species []models.Species) {
	speciesMap := map[string]uint{}
	for _, s := range species {
		speciesMap[s.Name] = s.ID
	}

	vaccines := []models.Vaccine{
		// Cattle
		{SpeciesID: speciesMap["Cattle"], Name: "BVD", Dose: s("2ml"), MinimumAgeValue: 2, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Cattle"], Name: "IBR", Dose: s("2ml"), MinimumAgeValue: 2, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Cattle"], Name: "Blackleg", Dose: s("5ml"), MinimumAgeValue: 3, MinimumAgeUnit: "Month", IntervalValue: 6, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Cattle"], Name: "Anthrax", Dose: s("1ml"), MinimumAgeValue: 4, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Cattle"], Name: "FMD", Dose: s("2ml"), MinimumAgeValue: 4, MinimumAgeUnit: "Month", IntervalValue: 6, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Cattle"], Name: "Brucellosis", Dose: s("2ml"), MinimumAgeValue: 4, MinimumAgeUnit: "Month", IntervalValue: 0, IntervalUnit: "Year", IsRepeatable: false},
		{SpeciesID: speciesMap["Cattle"], Name: "Rabies", Dose: s("1ml"), MinimumAgeValue: 3, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Cattle"], Name: "Clostridial 7-Way", Dose: s("5ml"), MinimumAgeValue: 3, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},

		// Goat
		{SpeciesID: speciesMap["Goat"], Name: "PPR", Dose: s("1ml"), MinimumAgeValue: 3, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Goat"], Name: "Goat Pox", Dose: s("0.5ml"), MinimumAgeValue: 2, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Goat"], Name: "Enterotoxemia", Dose: s("2ml"), MinimumAgeValue: 2, MinimumAgeUnit: "Month", IntervalValue: 6, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Goat"], Name: "Anthrax", Dose: s("1ml"), MinimumAgeValue: 4, MinimumAgeUnit: "Month", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Goat"], Name: "FMD", Dose: s("2ml"), MinimumAgeValue: 4, MinimumAgeUnit: "Month", IntervalValue: 6, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Goat"], Name: "Brucellosis", Dose: s("2ml"), MinimumAgeValue: 3, MinimumAgeUnit: "Month", IntervalValue: 0, IntervalUnit: "Year", IsRepeatable: false},

		// Chicken
		{SpeciesID: speciesMap["Chicken"], Name: "Newcastle Disease", Dose: s("0.5ml"), MinimumAgeValue: 1, MinimumAgeUnit: "Week", IntervalValue: 4, IntervalUnit: "Week", IsRepeatable: true},
		{SpeciesID: speciesMap["Chicken"], Name: "Gumboro (IBD)", Dose: s("0.5ml"), MinimumAgeValue: 2, MinimumAgeUnit: "Week", IntervalValue: 0, IntervalUnit: "Year", IsRepeatable: false},
		{SpeciesID: speciesMap["Chicken"], Name: "Fowl Pox", Dose: s("0.3ml"), MinimumAgeValue: 6, MinimumAgeUnit: "Week", IntervalValue: 12, IntervalUnit: "Month", IsRepeatable: true},
		{SpeciesID: speciesMap["Chicken"], Name: "Marek's Disease", Dose: s("0.2ml"), MinimumAgeValue: 1, MinimumAgeUnit: "Day", IntervalValue: 0, IntervalUnit: "Year", IsRepeatable: false},
		{SpeciesID: speciesMap["Chicken"], Name: "Infectious Bronchitis", Dose: s("0.5ml"), MinimumAgeValue: 1, MinimumAgeUnit: "Week", IntervalValue: 4, IntervalUnit: "Week", IsRepeatable: true},
		{SpeciesID: speciesMap["Chicken"], Name: "Avian Influenza", Dose: s("0.5ml"), MinimumAgeValue: 2, MinimumAgeUnit: "Week", IntervalValue: 6, IntervalUnit: "Month", IsRepeatable: true},
	}
	DB.Create(&vaccines)
	log.Printf("Seeded %d vaccines", len(vaccines))
}
