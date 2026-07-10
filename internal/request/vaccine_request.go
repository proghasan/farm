package request

type CreateVaccineRequest struct {
	SpeciesID       uint    `json:"species_id" validate:"required"`
	Name            string  `json:"name" validate:"required,min=1,max=150"`
	Description     *string `json:"description"`
	Dose            *string `json:"dose"`
	MinimumAgeValue uint    `json:"minimum_age_value" validate:"required"`
	MinimumAgeUnit  string  `json:"minimum_age_unit" validate:"required,oneof=Day Week Month Year"`
	IntervalValue   int     `json:"interval_value" validate:"required"`
	IntervalUnit    string  `json:"interval_unit" validate:"required,oneof=Day Week Month Year"`
	IsRepeatable    bool    `json:"is_repeatable"`
}

type UpdateVaccineRequest struct {
	Name            *string  `json:"name" validate:"omitempty,min=1,max=150"`
	Description     *string  `json:"description"`
	Dose            *string  `json:"dose"`
	MinimumAgeValue *uint    `json:"minimum_age_value"`
	MinimumAgeUnit  *string  `json:"minimum_age_unit" validate:"omitempty,oneof=Day Week Month Year"`
	IntervalValue   *int     `json:"interval_value"`
	IntervalUnit    *string  `json:"interval_unit" validate:"omitempty,oneof=Day Week Month Year"`
	IsRepeatable    *bool    `json:"is_repeatable"`
}
