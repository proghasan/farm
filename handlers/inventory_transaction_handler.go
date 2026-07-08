package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ListInventoryTransactions(c fiber.Ctx, db *gorm.DB) error {
	var txns []models.InventoryTransaction
	tx := db.Model(&models.InventoryTransaction{}).Preload("InventoryItem.Category")
	if itemID := c.Query("inventory_item_id"); itemID != "" {
		tx = tx.Where("inventory_item_id = ?", itemID)
	}
	if txnType := c.Query("transaction_type"); txnType != "" {
		tx = tx.Where("transaction_type = ?", txnType)
	}
	return paginate(c, tx.Order("transaction_date DESC"), &txns)
}

func GetInventoryTransaction(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var txn models.InventoryTransaction
	if err := db.Preload("InventoryItem.Category").First(&txn, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	return c.JSON(txn)
}

func CreateInventoryTransaction(c fiber.Ctx, db *gorm.DB) error {
	var txn models.InventoryTransaction
	if err := validateBody(c, &txn); err != nil {
		return err
	}
	txn.CreatedBy = middleware.GetUserID(c)
	txn.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&txn).Error; err != nil {
		return handleError(c, err)
	}
	db.Preload("InventoryItem.Category").First(&txn, txn.ID)
	return c.Status(201).JSON(txn)
}

func UpdateInventoryTransaction(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var txn models.InventoryTransaction
	if err := db.First(&txn, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	var input models.InventoryTransaction
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = txn.ID
	input.CreatedBy = txn.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&txn).Updates(input)
	db.Preload("InventoryItem.Category").First(&txn, txn.ID)
	return c.JSON(txn)
}

func DeleteInventoryTransaction(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	if err := db.Delete(&models.InventoryTransaction{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
