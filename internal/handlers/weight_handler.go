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

func NewWeightHandler(db *gorm.DB) *WeightHandler {
	return &WeightHandler{
		repo:       repositories.NewWeightRepository(db),
		animalRepo: repositories.NewAnimalRepository(db),
		db:         db,
	}
}

type WeightHandler struct {
	repo       *repositories.WeightRepository
	animalRepo *repositories.AnimalRepository
	db         *gorm.DB
}

func (h *WeightHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	animalID := c.Query("animal_id")
	weights, total, err := h.repo.List(animalID, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = weights
	return c.JSON(resp)
}

func (h *WeightHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	w, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	return c.JSON(w)
}

func (h *WeightHandler) Create(c fiber.Ctx) error {
	var req request.CreateWeightRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	w := models.AnimalWeightHistory{
		AnimalID:   req.AnimalID,
		Weight:     req.Weight,
		RecordDate: req.RecordDate,
		Remarks:    req.Remarks,
	}
	w.CreatedBy = middleware.GetUserID(c)
	w.UpdatedBy = middleware.GetUserID(c)

	err := h.db.Transaction(func(tx *gorm.DB) error {
		wRepo := repositories.NewWeightRepository(tx)
		aRepo := repositories.NewAnimalRepository(tx)
		if err := wRepo.Create(&w); err != nil {
			return err
		}
		aRepo.UpdateCurrentWeight(w.AnimalID, w.Weight)
		return nil
	})
	if err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&w)
	return c.Status(201).JSON(w)
}

func (h *WeightHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	w, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	var req request.UpdateWeightRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.RecordDate != nil {
		updates["record_date"] = *req.RecordDate
	}
	if req.Remarks != nil {
		updates["remarks"] = *req.Remarks
	}
	if err := h.repo.Update(w, updates); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.JSON(w)
}

func (h *WeightHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
