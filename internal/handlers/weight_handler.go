package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
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
	animalID := c.Query("animal_id")
	weights, err := h.repo.List(animalID)
	if err != nil {
		return err
	}
	return c.JSON(weights)
}

func (h *WeightHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	w, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	return c.JSON(w)
}

func (h *WeightHandler) Create(c fiber.Ctx) error {
	var req request.CreateWeightRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	w := models.AnimalWeightHistory{
		AnimalID:   req.AnimalID,
		Weight:     req.Weight,
		RecordDate: models.DateString(req.RecordDate),
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
		return err
	}
	h.repo.Preload(&w)
	return c.Status(201).JSON(w)
}

func (h *WeightHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	w, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	var req request.UpdateWeightRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
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
		return err
	}

	if req.Weight != nil {
		latest, err := h.repo.GetLatestByAnimalID(w.AnimalID)
		if err == nil && latest.ID == w.ID {
			h.animalRepo.UpdateCurrentWeight(w.AnimalID, *req.Weight)
		}
	}

	return c.JSON(w)
}

func (h *WeightHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)

	w, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}

	latest, err := h.repo.GetLatestByAnimalID(w.AnimalID)
	if err != nil {
		h.animalRepo.UpdateCurrentWeight(w.AnimalID, 0)
	} else {
		h.animalRepo.UpdateCurrentWeight(w.AnimalID, latest.Weight)
	}

	return c.SendStatus(204)
}
