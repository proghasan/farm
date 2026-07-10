package request

type CreateSpeciesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateSpeciesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}
