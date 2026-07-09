package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ListPregnancies(c fiber.Ctx, db *gorm.DB) error {
	var items []models.AnimalPregnancy
	tx := db.Model(&models.AnimalPregnancy{}).
		Preload("Animal").
		Preload("Breeder")
	if animalID := c.Query("animal_id"); animalID != "" {
		tx = tx.Where("animal_id = ?", animalID)
	}
	if status := c.Query("status"); status != "" {
		tx = tx.Where("status = ?", status)
	}
	return paginate(c, tx.Order("created_at DESC"), &items)
}

func GetPregnancy(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var item models.AnimalPregnancy
	if err := db.Preload("Animal").Preload("Breeder").First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	return c.JSON(item)
}

func CreatePregnancy(c fiber.Ctx, db *gorm.DB) error {
	var item models.AnimalPregnancy
	if err := validateBody(c, &item); err != nil {
		return err
	}
	item.CreatedBy = middleware.GetUserID(c)
	item.UpdatedBy = middleware.GetUserID(c)
	db.Create(&item)
	db.Preload("Animal").Preload("Breeder").First(&item, item.ID)
	return c.Status(201).JSON(item)
}

func UpdatePregnancy(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var item models.AnimalPregnancy
	if err := db.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	var input models.AnimalPregnancy
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = item.ID
	input.CreatedBy = item.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&item).Updates(input)
	db.Preload("Animal").Preload("Breeder").First(&item, item.ID)
	return c.JSON(item)
}

func DeletePregnancy(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	if err := db.Delete(&models.AnimalPregnancy{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}