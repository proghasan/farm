package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListSpecies(c *fiber.Ctx, db *gorm.DB) error {
	var species []models.Species
	tx := db.Model(&models.Species{}).Preload("Breeds")
	return paginate(c, tx, &species)
}

func GetSpecies(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var species models.Species
	if err := db.Preload("Breeds").First(&species, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
	}
	return c.JSON(species)
}

func CreateSpecies(c *fiber.Ctx, db *gorm.DB) error {
	var species models.Species
	if err := validateBody(c, &species); err != nil {
		return err
	}
	species.CreatedBy = middleware.GetUserID(c)
	species.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&species).Error; err != nil {
		return handleError(c, err)
	}
	return c.Status(201).JSON(species)
}

func UpdateSpecies(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var species models.Species
	if err := db.First(&species, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
	}
	var input models.Species
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = species.ID
	input.CreatedBy = species.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&species).Updates(input)
	return c.JSON(species)
}

func DeleteSpecies(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.Species{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
