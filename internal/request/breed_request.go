package request

type CreateBreedRequest struct {
	SpeciesID uint   `json:"species_id" validate:"required"`
	Name      string `json:"name" validate:"required,min=1,max=150"`
}

type UpdateBreedRequest struct {
	SpeciesID *uint   `json:"species_id"`
	Name      *string `json:"name" validate:"omitempty,min=1,max=150"`
}
