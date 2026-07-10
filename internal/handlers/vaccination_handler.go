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

func NewVaccinationHandler(db *gorm.DB) *VaccinationHandler {
	return &VaccinationHandler{repo: repositories.NewVaccinationRepository(db)}
}

type VaccinationHandler struct {
	repo *repositories.VaccinationRepository
}

func (h *VaccinationHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	animalID := c.Query("animal_id")
	vaccs, total, err := h.repo.List(animalID, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = vaccs
	return c.JSON(resp)
}

func (h *VaccinationHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	v, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccination not found"})
	}
	return c.JSON(v)
}

func (h *VaccinationHandler) Create(c fiber.Ctx) error {
	var v models.AnimalVaccination
	if err := validator.Body(c, &v); err != nil {
		return nil
	}
	v.CreatedBy = middleware.GetUserID(c)
	v.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&v); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&v)
	return c.Status(201).JSON(v)
}

func (h *VaccinationHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	v, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccination not found"})
	}
	var input models.AnimalVaccination
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.repo.Update(v, map[string]interface{}{
		"animal_id":       input.AnimalID,
		"vaccine_id":      input.VaccineID,
		"vaccination_date": input.VaccinationDate,
		"next_due_date":   input.NextDueDate,
		"doctor_name":     input.DoctorName,
		"remarks":         input.Remarks,
		"updated_by":      middleware.GetUserID(c),
	}); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(v)
	return c.JSON(v)
}

func (h *VaccinationHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
