package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
	"farm/internal/response"
	"farm/internal/validator"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{repo: repositories.NewUserRepository(db)}
}

type UserHandler struct {
	repo *repositories.UserRepository
}

func (h *UserHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	users, total, err := h.repo.List(page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = users
	return c.JSON(resp)
}

func (h *UserHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	user, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func (h *UserHandler) Profile(c fiber.Ctx) error {
	uid := middleware.GetUserID(c)
	user, err := h.repo.GetByIDFull(uid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
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
	if err := h.repo.Create(&user); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.Status(201).JSON(response.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	uid := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)

	if uid != uint(id) && role != "Owner" && role != "Manager" {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}

	user, err := h.repo.GetByIDFull(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var req request.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
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
	if req.Role != nil && role == "Owner" {
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

	if err := h.repo.Update(user, updates); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.JSON(response.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (h *UserHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
