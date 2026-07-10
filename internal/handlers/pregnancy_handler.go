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

func NewPregnancyHandler(db *gorm.DB) *PregnancyHandler {
	return &PregnancyHandler{repo: repositories.NewPregnancyRepository(db)}
}

type PregnancyHandler struct {
	repo *repositories.PregnancyRepository
}

func (h *PregnancyHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	animalID := c.Query("animal_id")
	status := c.Query("status")
	items, total, err := h.repo.List(animalID, status, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = items
	return c.JSON(resp)
}

func (h *PregnancyHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	item, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	return c.JSON(item)
}

func (h *PregnancyHandler) Create(c fiber.Ctx) error {
	var item models.AnimalPregnancy
	if err := validator.Body(c, &item); err != nil {
		return nil
	}
	item.CreatedBy = middleware.GetUserID(c)
	item.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&item); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&item)
	return c.Status(201).JSON(item)
}

func (h *PregnancyHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	item, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	var input models.AnimalPregnancy
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.repo.Update(item, map[string]interface{}{
		"animal_id":               input.AnimalID,
		"breeder_id":              input.BreederID,
		"mating_date":             input.MatingDate,
		"expected_due_date":       input.ExpectedDueDate,
		"actual_birth_date":       input.ActualBirthDate,
		"status":                  input.Status,
		"note":                    input.Note,
		"number_of_children":      input.NumberOfChildren,
		"number_of_male_children":  input.NumberOfMaleChildren,
		"number_of_female_children": input.NumberOfFemaleChildren,
		"number_of_dead_children":  input.NumberOfDeadChildren,
		"updated_by":              middleware.GetUserID(c),
	}); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(item)
	return c.JSON(item)
}

func (h *PregnancyHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
