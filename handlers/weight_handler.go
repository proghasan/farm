package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ListWeightHistories(c fiber.Ctx, db *gorm.DB) error {
	var weights []models.AnimalWeightHistory
	tx := db.Model(&models.AnimalWeightHistory{}).Preload("Animal")
	if animalID := c.Query("animal_id"); animalID != "" {
		tx = tx.Where("animal_id = ?", animalID)
	}
	return paginate(c, tx.Order("record_date DESC"), &weights)
}

func GetWeightHistory(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var w models.AnimalWeightHistory
	if err := db.Preload("Animal").First(&w, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	return c.JSON(w)
}

func CreateWeightHistory(c fiber.Ctx, db *gorm.DB) error {
	var w models.AnimalWeightHistory
	if err := validateBody(c, &w); err != nil {
		return nil
	}
	w.CreatedBy = middleware.GetUserID(c)
	w.UpdatedBy = middleware.GetUserID(c)

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&w).Error; err != nil {
			return err
		}
		tx.Model(&models.Animal{}).Where("id = ?", w.AnimalID).Update("current_weight", w.Weight)
		return nil
	})
	if err != nil {
		return handleError(c, err)
	}
	db.Preload("Animal").First(&w, w.ID)
	return c.Status(201).JSON(w)
}

func UpdateWeightHistory(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var w models.AnimalWeightHistory
	if err := db.First(&w, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	var input models.AnimalWeightHistory
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = w.ID
	input.CreatedBy = w.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&w).Updates(input)
	return c.JSON(w)
}

func DeleteWeightHistory(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	if err := db.Delete(&models.AnimalWeightHistory{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
