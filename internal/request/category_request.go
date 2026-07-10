package request

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=150"`
}

type UpdateCategoryRequest struct {
	Name *string `json:"name"`
}
