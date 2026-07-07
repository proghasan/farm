package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListInventoryItems(c *fiber.Ctx, db *gorm.DB) error {
	var items []models.InventoryItem
	tx := db.Model(&models.InventoryItem{}).Preload("Category")
	if categoryID := c.Query("category_id"); categoryID != "" {
		tx = tx.Where("category_id = ?", categoryID)
	}
	return paginate(c, tx, &items)
}

func GetInventoryItem(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var item models.InventoryItem
	if err := db.Preload("Category").First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}
	return c.JSON(item)
}

func CreateInventoryItem(c *fiber.Ctx, db *gorm.DB) error {
	var item models.InventoryItem
	if err := validateBody(c, &item); err != nil {
		return err
	}
	item.CreatedBy = middleware.GetUserID(c)
	item.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&item).Error; err != nil {
		return handleError(c, err)
	}
	db.Preload("Category").First(&item, item.ID)
	return c.Status(201).JSON(item)
}

func UpdateInventoryItem(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var item models.InventoryItem
	if err := db.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}
	var input models.InventoryItem
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = item.ID
	input.CreatedBy = item.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&item).Updates(input)
	db.Preload("Category").First(&item, item.ID)
	return c.JSON(item)
}

func DeleteInventoryItem(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.InventoryItem{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
