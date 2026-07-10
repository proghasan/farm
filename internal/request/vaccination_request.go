package request

type CreateVaccinationRequest struct {
	AnimalID       uint    `json:"animal_id" validate:"required"`
	VaccineID      uint    `json:"vaccine_id" validate:"required"`
	VaccinationDate string `json:"vaccination_date" validate:"required"`
	NextDueDate    *string `json:"next_due_date"`
	DoctorName     *string `json:"doctor_name"`
	Remarks        *string `json:"remarks"`
}

type UpdateVaccinationRequest struct {
	VaccinationDate *string `json:"vaccination_date"`
	NextDueDate     *string `json:"next_due_date"`
	DoctorName      *string `json:"doctor_name"`
	Remarks         *string `json:"remarks"`
}
