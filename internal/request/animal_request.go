package request

type CreateAnimalRequest struct {
	TagNo         string  `json:"tag_no" validate:"required,min=1,max=50"`
	SpeciesID     uint    `json:"species_id" validate:"required"`
	BreedID       *uint   `json:"breed_id"`
	FatherID      *uint   `json:"father_id"`
	MotherID      *uint   `json:"mother_id"`
	Gender        string  `json:"gender" validate:"required,oneof=Male Female"`
	BirthDate     *string `json:"birth_date"`
	PurchaseDate  *string `json:"purchase_date"`
	PurchasePrice float64 `json:"purchase_price"`
	CurrentWeight *float64 `json:"current_weight"`
	Color         *string `json:"color"`
	Status        string  `json:"status" validate:"omitempty,oneof=Healthy Pregnant Sick Sold Dead"`
	Remarks       *string `json:"remarks"`
}

type UpdateAnimalRequest struct {
	TagNo         *string  `json:"tag_no"`
	SpeciesID     *uint    `json:"species_id"`
	BreedID       *uint    `json:"breed_id"`
	FatherID      *uint    `json:"father_id"`
	MotherID      *uint    `json:"mother_id"`
	Gender        *string  `json:"gender"`
	BirthDate     *string  `json:"birth_date"`
	PurchaseDate  *string  `json:"purchase_date"`
	PurchasePrice *float64 `json:"purchase_price"`
	CurrentWeight *float64 `json:"current_weight"`
	Color         *string  `json:"color"`
	Status        *string  `json:"status"`
	Remarks       *string  `json:"remarks"`
}
