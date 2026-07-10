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
	var req request.CreatePregnancyRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	item := models.AnimalPregnancy{
		AnimalID:        req.AnimalID,
		BreederID:       req.BreederID,
		MatingDate:      req.MatingDate,
		ExpectedDueDate: req.ExpectedDueDate,
		ActualBirthDate: req.ActualBirthDate,
		Status:          models.PregnancyStatus(req.Status),
		Note:            req.Note,
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
	var req request.UpdatePregnancyRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.BreederID != nil {
		updates["breeder_id"] = *req.BreederID
	}
	if req.MatingDate != nil {
		updates["mating_date"] = *req.MatingDate
	}
	if req.ExpectedDueDate != nil {
		updates["expected_due_date"] = *req.ExpectedDueDate
	}
	if req.ActualBirthDate != nil {
		updates["actual_birth_date"] = *req.ActualBirthDate
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.NumberOfChildren != nil {
		updates["number_of_children"] = *req.NumberOfChildren
	}
	if req.NumberOfMaleChildren != nil {
		updates["number_of_male_children"] = *req.NumberOfMaleChildren
	}
	if req.NumberOfFemaleChildren != nil {
		updates["number_of_female_children"] = *req.NumberOfFemaleChildren
	}
	if req.NumberOfDeadChildren != nil {
		updates["number_of_dead_children"] = *req.NumberOfDeadChildren
	}
	if err := h.repo.Update(item, updates); err != nil {
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
