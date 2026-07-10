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
	var req request.CreateVaccineRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	vaccine := models.Vaccine{
		SpeciesID:       req.SpeciesID,
		Name:            req.Name,
		Description:     req.Description,
		Dose:            req.Dose,
		MinimumAgeValue: req.MinimumAgeValue,
		MinimumAgeUnit:  req.MinimumAgeUnit,
		IntervalValue:   req.IntervalValue,
		IntervalUnit:    req.IntervalUnit,
		IsRepeatable:    req.IsRepeatable,
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
	var req request.UpdateVaccineRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Dose != nil {
		updates["dose"] = *req.Dose
	}
	if req.MinimumAgeValue != nil {
		updates["minimum_age_value"] = *req.MinimumAgeValue
	}
	if req.MinimumAgeUnit != nil {
		updates["minimum_age_unit"] = *req.MinimumAgeUnit
	}
	if req.IntervalValue != nil {
		updates["interval_value"] = *req.IntervalValue
	}
	if req.IntervalUnit != nil {
		updates["interval_unit"] = *req.IntervalUnit
	}
	if req.IsRepeatable != nil {
		updates["is_repeatable"] = *req.IsRepeatable
	}
	if err := h.repo.Update(vaccine, updates); err != nil {
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
