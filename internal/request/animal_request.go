package request

import "strconv"

type AnimalRequest struct {
	TagNo         string   `json:"tag_no"`
	SpeciesID     uint     `json:"species_id"`
	BreedID       uint   `json:"breed_id"`
	FatherID      *uint    `json:"father_id"`
	MotherID      *uint    `json:"mother_id"`
	Gender        string   `json:"gender"`
	BirthDate     *string  `json:"birth_date"`
	PurchaseDate  *string  `json:"purchase_date"`
	PurchasePrice float64  `json:"purchase_price"`
	CurrentWeight *float64 `json:"current_weight"`
	LastVaccine   *string  `json:"last_vaccine"`
	Color         *string  `json:"color"`
	Status        string   `json:"status"`
	Remarks       *string  `json:"remarks"`
}

func AnimalCreateRules() Rules {
	return Rules{
		"tag_no":      "required|min:1|max:50|unique:animals,tag_no",
		"species_id":  "required",
		"breed_id":    "required",
		"gender":      "required|min:4|max:6",
		"status":      "required",
	}
}

func AnimalUpdateRules(id int) Rules {
	return Rules{
		"tag_no":      "required|min:1|max:50|unique:animals,tag_no," + strconv.Itoa(id),
		"species_id":  "required",
		"breed_id":    "required",
		"gender":      "required|min:4|max:6",
		"status":      "required",
	}
}
