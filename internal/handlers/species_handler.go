package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
	"farm/internal/response"
	"farm/internal/validator"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewSpeciesHandler(db *gorm.DB) *SpeciesHandler {
	return &SpeciesHandler{repo: repositories.NewSpeciesRepository(db)}
}

type SpeciesHandler struct {
	repo *repositories.SpeciesRepository
}

func (h *SpeciesHandler) List(c fiber.Ctx) error {
	search := c.Query("search")
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	species, total, err := h.repo.List(search, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = species
	return c.JSON(resp)
}

func (h *SpeciesHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	species, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
	}
	return c.JSON(species)
}

func (h *SpeciesHandler) Create(c fiber.Ctx) error {
	var req request.CreateSpeciesRequest

	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}

	species := models.Species{
		Name: req.Name,
	}
	species.CreatedBy = middleware.GetUserID(c)
	species.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&species); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.Status(201).JSON(species)
}

func (h *SpeciesHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	species, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
	}
	
	var req request.UpdateSpeciesRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}

	species.Name = req.Name
	species.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Update(species); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.JSON(species)
}

func (h *SpeciesHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
