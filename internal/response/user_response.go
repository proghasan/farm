package response

type User struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
	Role     string  `json:"role"`
}
