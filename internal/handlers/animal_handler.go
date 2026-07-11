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
	animals, total, err := h.repo.List(filters, page, perPage)
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
	var req request.CreateAnimalRequest
	if err := c.Bind().Body(&req); err != nil {
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
		Color:         req.Color,
		Status:        models.AnimalStatus(req.Status),
		Remarks:       req.Remarks,
	}
	animal.CreatedBy = middleware.GetUserID(c)
	animal.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&animal); err != nil {
		return err
	}
	h.repo.Preload(&animal)
	return c.Status(201).JSON(animal)
}

func (h *AnimalHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	animal, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	var req request.UpdateAnimalRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.TagNo != nil {
		updates["tag_no"] = *req.TagNo
	}
	if req.SpeciesID != nil {
		updates["species_id"] = *req.SpeciesID
	}
	if req.BreedID != nil {
		updates["breed_id"] = *req.BreedID
	}
	if req.FatherID != nil {
		updates["father_id"] = *req.FatherID
	}
	if req.MotherID != nil {
		updates["mother_id"] = *req.MotherID
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.BirthDate != nil {
		updates["birth_date"] = *req.BirthDate
	}
	if req.PurchaseDate != nil {
		updates["purchase_date"] = *req.PurchaseDate
	}
	if req.PurchasePrice != nil {
		updates["purchase_price"] = *req.PurchasePrice
	}
	if req.CurrentWeight != nil {
		updates["current_weight"] = *req.CurrentWeight
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remarks != nil {
		updates["remarks"] = *req.Remarks
	}
	if err := h.repo.Update(animal, updates); err != nil {
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
	return c.SendStatus(204)
}
