package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
	"farm/internal/response"
	"farm/pkg/validator"

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
		return err
	}
	return c.JSON(species)
}

func (h *SpeciesHandler) Create(c fiber.Ctx) error {
	form := validator.New(c).Rules(request.SpeciesCreateRules())

	var req request.SpeciesRequest

	if err := form.Validate(&req); err != nil {
		return err
	}

	species := models.Species{
		Name:      req.Name,
		CreatedBy: middleware.GetUserID(c),
		UpdatedBy: middleware.GetUserID(c),
	}

	if err := h.repo.Create(&species); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(species)
}

func (h *SpeciesHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	species, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}

	form := validator.New(c).Rules(request.SpeciesUpdateRules())
	
	var req request.SpeciesRequest

	if err := form.Validate(&req); err != nil {
		return err
	}	

	species.Name = req.Name
	species.UpdatedBy = middleware.GetUserID(c)

	if err := h.repo.Update(species); err != nil {
		return err
	}

	return c.JSON(species)
}

func (h *SpeciesHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)

	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}