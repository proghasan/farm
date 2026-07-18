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

func NewAnimalHandler(db *gorm.DB) *AnimalHandler {
	return &AnimalHandler{
		repo:     repositories.NewAnimalRepository(db),
		pregRepo: repositories.NewPregnancyRepository(db),
	}
}

type AnimalHandler struct {
	repo     *repositories.AnimalRepository
	pregRepo *repositories.PregnancyRepository
}

func (h *AnimalHandler) List(c fiber.Ctx) error {
	search := c.Query("search")
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	filters := map[string]string{
		"species_id": c.Query("species_id"),
		"status":     c.Query("status"),
		"gender":     c.Query("gender"),
	}
	animals, total, err := h.repo.List(search, filters, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = animals
	return c.JSON(resp)
}

func (h *AnimalHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	animal, err := h.repo.GetByID(uint(id))

	if err != nil {
		return err
	}
	return c.JSON(animal)
}

func (h *AnimalHandler) Profile(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	animal, err := h.repo.GetProfile(uint(id))
	if err != nil {
		return err
	}
	pregnancies, _ := h.pregRepo.ListByAnimalID(uint(id))
	return c.JSON(fiber.Map{
		"animal":      animal,
		"pregnancies": pregnancies,
	})
}

func (h *AnimalHandler) Create(c fiber.Ctx) error {
	var req request.AnimalRequest

	if err := validator.New(c, h.repo.DB).Rules(request.AnimalCreateRules()).Validate(&req); err != nil {
		return err
	}

	animal := models.Animal{
		TagNo:         req.TagNo,
		SpeciesID:     req.SpeciesID,
		BreedID:       req.BreedID,
		FatherID:      req.FatherID,
		MotherID:      req.MotherID,
		Gender:        req.Gender,
		BirthDate:     req.BirthDate,
		PurchaseDate:  req.PurchaseDate,
		PurchasePrice: req.PurchasePrice,
		CurrentWeight: req.CurrentWeight,
		LastVaccine:   req.LastVaccine,
		Color:         req.Color,
		Status:        models.AnimalStatus(req.Status),
		Remarks:       req.Remarks,
		CreatedBy:     middleware.GetUserID(c),
		UpdatedBy:     middleware.GetUserID(c),
	}

	if err := h.repo.Create(&animal); err != nil {
		return err
	}

	h.repo.Preload(&animal)
	return c.Status(fiber.StatusCreated).JSON(animal)
}

func (h *AnimalHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	animal, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}

	var req request.AnimalRequest

	if err := validator.New(c, h.repo.DB).Rules(request.AnimalUpdateRules(id)).Validate(&req); err != nil {
		return err
	}

	animal.TagNo = req.TagNo
	animal.SpeciesID = req.SpeciesID
	animal.BreedID = req.BreedID
	animal.FatherID = req.FatherID
	animal.MotherID = req.MotherID
	animal.Gender = req.Gender
	animal.BirthDate = req.BirthDate
	animal.PurchaseDate = req.PurchaseDate
	animal.PurchasePrice = req.PurchasePrice
	animal.CurrentWeight = req.CurrentWeight
	animal.LastVaccine = req.LastVaccine
	animal.Color = req.Color
	animal.Status = models.AnimalStatus(req.Status)
	animal.Remarks = req.Remarks
	animal.UpdatedBy = middleware.GetUserID(c)
	animal.Species = nil
	animal.Breed = nil
	animal.User = models.User{}

	if err := h.repo.Update(animal); err != nil {
		return err
	}

	h.repo.Preload(animal)
	return c.JSON(animal)
}

func (h *AnimalHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)

	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
