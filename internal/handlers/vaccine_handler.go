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

func NewVaccineHandler(db *gorm.DB) *VaccineHandler {
	return &VaccineHandler{repo: repositories.NewVaccineRepository(db)}
}

type VaccineHandler struct {
	repo *repositories.VaccineRepository
}

func (h *VaccineHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	speciesID := c.Query("species_id")
	vaccines, total, err := h.repo.List(speciesID, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = vaccines
	return c.JSON(resp)
}

func (h *VaccineHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	vaccine, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccine not found"})
	}
	return c.JSON(vaccine)
}

func (h *VaccineHandler) Create(c fiber.Ctx) error {
	var vaccine models.Vaccine
	if err := validator.Body(c, &vaccine); err != nil {
		return nil
	}
	vaccine.CreatedBy = middleware.GetUserID(c)
	vaccine.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&vaccine); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&vaccine)
	return c.Status(201).JSON(vaccine)
}

func (h *VaccineHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	vaccine, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccine not found"})
	}
	var input models.Vaccine
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.repo.Update(vaccine, map[string]interface{}{
		"species_id":       input.SpeciesID,
		"name":             input.Name,
		"description":      input.Description,
		"dose":             input.Dose,
		"minimum_age_value": input.MinimumAgeValue,
		"minimum_age_unit": input.MinimumAgeUnit,
		"interval_value":   input.IntervalValue,
		"interval_unit":    input.IntervalUnit,
		"is_repeatable":    input.IsRepeatable,
		"updated_by":       middleware.GetUserID(c),
	}); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(vaccine)
	return c.JSON(vaccine)
}

func (h *VaccineHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
