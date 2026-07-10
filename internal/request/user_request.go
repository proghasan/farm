package request

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=150"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
	Password string  `json:"password" validate:"required,min=6"`
	Role     string  `json:"role" validate:"omitempty,oneof=Owner Manager Veterinarian Worker Accountant"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
	Status   *string `json:"status"`
}


