package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListVaccines(c *fiber.Ctx, db *gorm.DB) error {
	var vaccines []models.Vaccine
	tx := db.Model(&models.Vaccine{}).Preload("Species")
	if speciesID := c.Query("species_id"); speciesID != "" {
		tx = tx.Where("species_id = ?", speciesID)
	}
	return paginate(c, tx, &vaccines)
}

func GetVaccine(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var vaccine models.Vaccine
	if err := db.Preload("Species").First(&vaccine, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccine not found"})
	}
	return c.JSON(vaccine)
}

func CreateVaccine(c *fiber.Ctx, db *gorm.DB) error {
	var vaccine models.Vaccine
	if err := c.BodyParser(&vaccine); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	vaccine.CreatedBy = middleware.GetUserID(c)
	vaccine.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&vaccine).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(vaccine)
}

func UpdateVaccine(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var vaccine models.Vaccine
	if err := db.First(&vaccine, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vaccine not found"})
	}
	var input models.Vaccine
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = vaccine.ID
	input.CreatedBy = vaccine.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&vaccine).Updates(input)
	return c.JSON(vaccine)
}

func DeleteVaccine(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.Vaccine{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
