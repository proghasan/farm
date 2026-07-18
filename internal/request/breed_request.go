package request

import "strconv"

type BreedRequest struct {
	Name      string `json:"name"`
	SpeciesID uint   `json:"species_id"`
}

func BreedCreateRules() Rules {
	return Rules{
		"name":       "required|min:1|max:150|unique:breeds,name",
		"species_id": "required",
	}
}

func BreedUpdateRules(id int) Rules {
	return Rules{
		"name":       "required|min:1|max:150|unique:breeds,name," + strconv.Itoa(id),
		"species_id": "required",
	}
}
