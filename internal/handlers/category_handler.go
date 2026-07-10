package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/response"
	"farm/internal/validator"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{repo: repositories.NewCategoryRepository(db)}
}

type CategoryHandler struct {
	repo *repositories.CategoryRepository
}

func (h *CategoryHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	cats, total, err := h.repo.List(page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = cats
	return c.JSON(resp)
}

func (h *CategoryHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	cat, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(cat)
}

func (h *CategoryHandler) Create(c fiber.Ctx) error {
	var cat models.InventoryCategory
	if err := validator.Body(c, &cat); err != nil {
		return nil
	}
	cat.CreatedBy = middleware.GetUserID(c)
	cat.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&cat); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.Status(201).JSON(cat)
}

func (h *CategoryHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	cat, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	var input models.InventoryCategory
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.repo.Update(cat, map[string]interface{}{
		"name":       input.Name,
		"updated_by": middleware.GetUserID(c),
	}); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.JSON(cat)
}

func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
