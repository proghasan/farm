package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ListBreeds(c fiber.Ctx, db *gorm.DB) error {
	var breeds []models.Breed
	tx := db.Model(&models.Breed{}).Preload("Species")
	if speciesID := c.Query("species_id"); speciesID != "" {
		tx = tx.Where("species_id = ?", speciesID)
	}
	return paginate(c, tx, &breeds)
}

func GetBreed(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var breed models.Breed
	if err := db.Preload("Species").First(&breed, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Breed not found"})
	}
	return c.JSON(breed)
}

func CreateBreed(c fiber.Ctx, db *gorm.DB) error {
	var breed models.Breed
	if err := validateBody(c, &breed); err != nil {
		return nil
	}
	breed.CreatedBy = middleware.GetUserID(c)
	breed.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&breed).Error; err != nil {
		return handleError(c, err)
	}
	db.Preload("Species").First(&breed, breed.ID)
	return c.Status(201).JSON(breed)
}

func UpdateBreed(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var breed models.Breed
	if err := db.First(&breed, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Breed not found"})
	}
	var input models.Breed
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = breed.ID
	input.CreatedBy = breed.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&breed).Updates(input)
	db.Preload("Species").First(&breed, breed.ID)
	return c.JSON(breed)
}

func DeleteBreed(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	if err := db.Delete(&models.Breed{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
