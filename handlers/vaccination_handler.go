package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListVaccinations(c *fiber.Ctx, db *gorm.DB) error {
	var vaccinations []models.AnimalVaccination
	tx := db.Model(&models.AnimalVaccination{}).Preload("Animal").Preload("Vaccine")
	if animalID := c.Query("animal_id"); animalID != "" {
		tx = tx.Where("animal_id = ?", animalID)
	}
	return paginate(c, tx.Order("vaccination_date DESC"), &vaccinations)
}

func GetVaccination(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var v models.AnimalVaccination
	if err := db.Preload("Animal").Preload("Vaccine").First(&v, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccination not found"})
	}
	return c.JSON(v)
}

func CreateVaccination(c *fiber.Ctx, db *gorm.DB) error {
	var v models.AnimalVaccination
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	v.CreatedBy = middleware.GetUserID(c)
	v.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&v).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(v)
}

func UpdateVaccination(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var v models.AnimalVaccination
	if err := db.First(&v, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccination not found"})
	}
	var input models.AnimalVaccination
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = v.ID
	input.CreatedBy = v.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&v).Updates(input)
	return c.JSON(v)
}

func DeleteVaccination(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.AnimalVaccination{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
