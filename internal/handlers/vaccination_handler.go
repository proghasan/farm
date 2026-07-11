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
		return err
	}
	return c.JSON(v)
}

func (h *VaccinationHandler) Create(c fiber.Ctx) error {
	var req request.CreateVaccinationRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	v := models.AnimalVaccination{
		AnimalID:       req.AnimalID,
		VaccineID:      req.VaccineID,
		VaccinationDate: req.VaccinationDate,
		NextDueDate:    req.NextDueDate,
		DoctorName:     req.DoctorName,
		Remarks:        req.Remarks,
	}
	v.CreatedBy = middleware.GetUserID(c)
	v.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&v); err != nil {
		return err
	}
	h.repo.Preload(&v)
	return c.Status(201).JSON(v)
}

func (h *VaccinationHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	v, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	var req request.UpdateVaccinationRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.VaccinationDate != nil {
		updates["vaccination_date"] = *req.VaccinationDate
	}
	if req.NextDueDate != nil {
		updates["next_due_date"] = *req.NextDueDate
	}
	if req.DoctorName != nil {
		updates["doctor_name"] = *req.DoctorName
	}
	if req.Remarks != nil {
		updates["remarks"] = *req.Remarks
	}
	if err := h.repo.Update(v, updates); err != nil {
		return err
	}
	h.repo.Preload(v)
	return c.JSON(v)
}

func (h *VaccinationHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(204)
}
