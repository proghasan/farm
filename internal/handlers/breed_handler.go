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

func NewBreedHandler(db *gorm.DB) *BreedHandler {
	return &BreedHandler{repo: repositories.NewBreedRepository(db)}
}

type BreedHandler struct {
	repo *repositories.BreedRepository
}

func (h *BreedHandler) List(c fiber.Ctx) error {
	if c.Query("all") == "true" {
		breeds, err := h.repo.All()
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"data": breeds})
	}

	search := c.Query("search")
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	breeds, total, err := h.repo.List(search, page, perPage)
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
	var req request.BreedRequest

	if err := validator.New(c, h.repo.DB).Rules(request.BreedCreateRules()).Validate(&req); err != nil {
		return err
	}

	breed := models.Breed{
		SpeciesID: req.SpeciesID,
		Name:      req.Name,
		CreatedBy: middleware.GetUserID(c),
		UpdatedBy: middleware.GetUserID(c),
	}

	if err := h.repo.Create(&breed); err != nil {
		return err
	}

	h.repo.Preload(&breed)
	return c.Status(fiber.StatusCreated).JSON(breed)
}

func (h *BreedHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	breed, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}

	var req request.BreedRequest

	if err := validator.New(c, h.repo.DB).Rules(request.BreedUpdateRules(id)).Validate(&req); err != nil {
		return err
	}

	breed.Name = req.Name
	breed.SpeciesID = req.SpeciesID
	breed.UpdatedBy = middleware.GetUserID(c)
	breed.Species = models.Species{}
	breed.User = models.User{}

	if err := h.repo.Update(breed); err != nil {
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

	return c.SendStatus(fiber.StatusNoContent)
}
