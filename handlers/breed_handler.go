package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListBreeds(c *fiber.Ctx, db *gorm.DB) error {
	var breeds []models.Breed
	tx := db.Preload("Species")
	if speciesID := c.Query("species_id"); speciesID != "" {
		tx = tx.Where("species_id = ?", speciesID)
	}
	return paginate(c, tx, &breeds)
}

func GetBreed(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var breed models.Breed
	if err := db.Preload("Species").First(&breed, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Breed not found"})
	}
	return c.JSON(breed)
}

func CreateBreed(c *fiber.Ctx, db *gorm.DB) error {
	var breed models.Breed
	if err := c.BodyParser(&breed); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	breed.CreatedBy = middleware.GetUserID(c)
	breed.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&breed).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(breed)
}

func UpdateBreed(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var breed models.Breed
	if err := db.First(&breed, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Breed not found"})
	}
	var input models.Breed
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = breed.ID
	input.CreatedBy = breed.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&breed).Updates(input)
	return c.JSON(breed)
}

func DeleteBreed(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.Breed{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
