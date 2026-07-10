package request

type CreateWeightRequest struct {
	AnimalID   uint    `json:"animal_id" validate:"required"`
	Weight     float64 `json:"weight" validate:"required,gt=0"`
	RecordDate string  `json:"record_date" validate:"required"`
	Remarks    *string `json:"remarks"`
}

type UpdateWeightRequest struct {
	Weight     *float64 `json:"weight"`
	RecordDate *string  `json:"record_date"`
	Remarks    *string  `json:"remarks"`
}
