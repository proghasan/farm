package handlers

import (
	"farm/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func GetAnimalProfile(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)

	var animal models.Animal
	if err := db.
		Preload("Species").
		Preload("Breed").
		Preload("Father").
		Preload("Mother").
		Preload("WeightHistories").
		Preload("AnimalVaccinations.Vaccine").
		First(&animal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
	}

	var pregnancies []models.AnimalPregnancy
	db.Where("animal_id = ?", id).Order("created_at DESC").Find(&pregnancies)

	return c.JSON(fiber.Map{
		"animal":      animal,
		"pregnancies": pregnancies,
	})
}