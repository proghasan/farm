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
		return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
	}
	return c.JSON(animal)
}

func (h *AnimalHandler) Profile(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	animal, err := h.repo.GetProfile(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
	}
	pregnancies, _ := h.pregRepo.ListByAnimalID(uint(id))
	return c.JSON(fiber.Map{
		"animal":      animal,
		"pregnancies": pregnancies,
	})
}

func (h *AnimalHandler) Create(c fiber.Ctx) error {
	var animal models.Animal
	if err := validator.Body(c, &animal); err != nil {
		return nil
	}
	animal.CreatedBy = middleware.GetUserID(c)
	animal.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&animal); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&animal)
	return c.Status(201).JSON(animal)
}

func (h *AnimalHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	animal, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
	}
	var input models.Animal
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.repo.Update(animal, map[string]interface{}{
		"tag_no":         input.TagNo,
		"species_id":     input.SpeciesID,
		"breed_id":       input.BreedID,
		"father_id":      input.FatherID,
		"mother_id":      input.MotherID,
		"gender":         input.Gender,
		"birth_date":     input.BirthDate,
		"purchase_date":  input.PurchaseDate,
		"purchase_price": input.PurchasePrice,
		"current_weight": input.CurrentWeight,
		"color":          input.Color,
		"status":         input.Status,
		"remarks":        input.Remarks,
		"updated_by":     middleware.GetUserID(c),
	}); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(animal)
	return c.JSON(animal)
}

func (h *AnimalHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
