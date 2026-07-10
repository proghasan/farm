package request

import "strings"

type CreateSpeciesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateSpeciesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (r *CreateSpeciesRequest) Trim() {
	r.Name = strings.TrimSpace(r.Name)
}

func (r *UpdateSpeciesRequest) Trim() {
	r.Name = strings.TrimSpace(r.Name)
}
