package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListInventoryCategories(c *fiber.Ctx, db *gorm.DB) error {
	var categories []models.InventoryCategory
	tx := db.Preload("Items")
	return paginate(c, tx, &categories)
}

func GetInventoryCategory(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var cat models.InventoryCategory
	if err := db.Preload("Items").First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(cat)
}

func CreateInventoryCategory(c *fiber.Ctx, db *gorm.DB) error {
	var cat models.InventoryCategory
	if err := c.BodyParser(&cat); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	cat.CreatedBy = middleware.GetUserID(c)
	cat.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&cat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(cat)
}

func UpdateInventoryCategory(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var cat models.InventoryCategory
	if err := db.First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	var input models.InventoryCategory
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = cat.ID
	input.CreatedBy = cat.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&cat).Updates(input)
	return c.JSON(cat)
}

func DeleteInventoryCategory(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.InventoryCategory{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
