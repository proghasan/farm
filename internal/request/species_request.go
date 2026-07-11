package request

import "strconv"

type SpeciesRequest struct {
	Name string
}

type Rules map[string]string

func SpeciesCreateRules() Rules {
	return Rules{
		"name": "required|unique:species,name|min:3|max:255",
	}
}

func SpeciesUpdateRules(id int) Rules {
	return Rules{
		"name": "required|unique:species,name," + strconv.Itoa(id) + "|min:3|max:255",
	}
}
