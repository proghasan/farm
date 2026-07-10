package request

type CreateAccountHeadRequest struct {
	Type        string  `json:"type" validate:"required,oneof=Income Expense"`
	Name        string  `json:"name" validate:"required,min=1,max=150"`
	Description *string `json:"description"`
}

type UpdateAccountHeadRequest struct {
	Type        *string `json:"type"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
