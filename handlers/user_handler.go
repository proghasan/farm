package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

func ListUsers(c fiber.Ctx, db *gorm.DB) error {
	var users []models.User
	tx := db.Model(&models.User{}).Select("id, name, email, phone, username, role, status, created_at, updated_at")
	return paginate(c, tx, &users)
}

func GetUser(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var user models.User
	if err := db.Select("id, name, email, phone, username, role, status, created_at, updated_at").First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func CreateUser(c fiber.Ctx, db *gorm.DB) error {
	var req CreateUserRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Username: req.Username,
		Password: string(hashed),
		Role:     req.Role,
	}
	if user.Role == "" {
		user.Role = "Worker"
	}

	if err := db.Create(&user).Error; err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"phone":    user.Phone,
		"username": user.Username,
		"role":     user.Role,
	})
}

func UpdateUser(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)

	// Only Owner and Manager can update other users; users can update themselves
	uid := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)
	if uid != uint(id) && role != "Owner" && role != "Manager" {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Role != nil && (role == "Owner") {
		updates["role"] = *req.Role
	}
	if req.Status != nil && (role == "Owner" || role == "Manager") {
		updates["status"] = *req.Status
	}
	if req.Password != nil && *req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		updates["password"] = string(hashed)
	}

	db.Model(&user).Updates(updates)
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"phone":    user.Phone,
		"username": user.Username,
		"role":     user.Role,
	})
}

func DeleteUser(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}

func GetProfile(c fiber.Ctx, db *gorm.DB) error {
	uid := middleware.GetUserID(c)
	var user models.User
	if err := db.Select("id, name, email, phone, username, role, status, avatar, last_login_at, created_at, updated_at").First(&user, uid).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}
