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

func NewBreedHandler(db *gorm.DB) *BreedHandler {
	return &BreedHandler{repo: repositories.NewBreedRepository(db)}
}

type BreedHandler struct {
	repo *repositories.BreedRepository
}

func (h *BreedHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	speciesID := c.Query("species_id")
	breeds, total, err := h.repo.List(speciesID, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = breeds
	return c.JSON(resp)
}

func (h *BreedHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	breed, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	return c.JSON(breed)
}

func (h *BreedHandler) Create(c fiber.Ctx) error {
	var req request.CreateBreedRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	breed := models.Breed{
		SpeciesID: req.SpeciesID,
		Name:      req.Name,
	}
	breed.CreatedBy = middleware.GetUserID(c)
	breed.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&breed); err != nil {
		return err
	}
	h.repo.Preload(&breed)
	return c.Status(201).JSON(breed)
}

func (h *BreedHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	breed, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	var req request.UpdateBreedRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.SpeciesID != nil {
		updates["species_id"] = *req.SpeciesID
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if err := h.repo.Update(breed, updates); err != nil {
		return err
	}
	h.repo.Preload(breed)
	return c.JSON(breed)
}

func (h *BreedHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(204)
}
