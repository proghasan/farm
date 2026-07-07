package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListInventoryTransactions(c *fiber.Ctx, db *gorm.DB) error {
	var txns []models.InventoryTransaction
	tx := db.Preload("InventoryItem.Category")
	if itemID := c.Query("inventory_item_id"); itemID != "" {
		tx = tx.Where("inventory_item_id = ?", itemID)
	}
	if txnType := c.Query("transaction_type"); txnType != "" {
		tx = tx.Where("transaction_type = ?", txnType)
	}
	return paginate(c, tx.Order("transaction_date DESC"), &txns)
}

func GetInventoryTransaction(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var txn models.InventoryTransaction
	if err := db.Preload("InventoryItem.Category").First(&txn, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	return c.JSON(txn)
}

func CreateInventoryTransaction(c *fiber.Ctx, db *gorm.DB) error {
	var txn models.InventoryTransaction
	if err := c.BodyParser(&txn); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	txn.CreatedBy = middleware.GetUserID(c)
	txn.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&txn).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(txn)
}

func UpdateInventoryTransaction(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var txn models.InventoryTransaction
	if err := db.First(&txn, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	var input models.InventoryTransaction
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = txn.ID
	input.CreatedBy = txn.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&txn).Updates(input)
	return c.JSON(txn)
}

func DeleteInventoryTransaction(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.InventoryTransaction{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
