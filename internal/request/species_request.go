package request

type SpeciesRequest struct {
	Name string
}

type Rules map[string]string

func SpeciesCreateRules() Rules {
	return Rules{
		"name": "required|unique|min:3|max:255",
	}
}

func SpeciesUpdateRules() Rules {
	return Rules{
		"name": "required|min:3|max:255",
	}
}
