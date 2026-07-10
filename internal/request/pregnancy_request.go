package request

type CreatePregnancyRequest struct {
	AnimalID        uint    `json:"animal_id" validate:"required"`
	BreederID       *uint   `json:"breeder_id"`
	MatingDate      string  `json:"mating_date" validate:"required"`
	ExpectedDueDate string  `json:"expected_due_date" validate:"required"`
	ActualBirthDate *string `json:"actual_birth_date"`
	Status          string  `json:"status" validate:"omitempty,oneof=Mated Pregnant Delivered Aborted Miscarriage Failed"`
	Note            *string `json:"note"`
}

type UpdatePregnancyRequest struct {
	BreederID             *uint   `json:"breeder_id"`
	MatingDate            *string `json:"mating_date"`
	ExpectedDueDate       *string `json:"expected_due_date"`
	ActualBirthDate       *string `json:"actual_birth_date"`
	Status                *string `json:"status" validate:"omitempty,oneof=Mated Pregnant Delivered Aborted Miscarriage Failed"`
	Note                  *string `json:"note"`
	NumberOfChildren      *int    `json:"number_of_children"`
	NumberOfMaleChildren  *int    `json:"number_of_male_children"`
	NumberOfFemaleChildren *int   `json:"number_of_female_children"`
	NumberOfDeadChildren  *int    `json:"number_of_dead_children"`
}
