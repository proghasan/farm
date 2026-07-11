package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
	"farm/internal/response"
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
		return err
	}
	return c.JSON(cat)
}

func (h *CategoryHandler) Create(c fiber.Ctx) error {
	var req request.CreateCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	cat := models.InventoryCategory{
		Name: req.Name,
	}
	cat.CreatedBy = middleware.GetUserID(c)
	cat.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&cat); err != nil {
		return err
	}
	return c.Status(201).JSON(cat)
}

func (h *CategoryHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	cat, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	var req request.UpdateCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if err := h.repo.Update(cat, updates); err != nil {
		return err
	}
	return c.JSON(cat)
}

func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(204)
}
