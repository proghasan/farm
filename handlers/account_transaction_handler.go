package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListAccountTransactions(c *fiber.Ctx, db *gorm.DB) error {
	var txns []models.AccountTransaction
	tx := db.Model(&models.AccountTransaction{}).Preload("AccountHead")
	if headID := c.Query("account_head_id"); headID != "" {
		tx = tx.Where("account_head_id = ?", headID)
	}
	if acctType := c.Query("type"); acctType != "" {
		tx = tx.Joins("AccountHead").Where("AccountHead.type = ?", acctType)
	}
	return paginate(c, tx.Order("transaction_date DESC"), &txns)
}

func GetAccountTransaction(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var txn models.AccountTransaction
	if err := db.Preload("AccountHead").First(&txn, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	return c.JSON(txn)
}

func CreateAccountTransaction(c *fiber.Ctx, db *gorm.DB) error {
	var txn models.AccountTransaction
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

func UpdateAccountTransaction(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var txn models.AccountTransaction
	if err := db.First(&txn, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	var input models.AccountTransaction
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = txn.ID
	input.CreatedBy = txn.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&txn).Updates(input)
	return c.JSON(txn)
}

func DeleteAccountTransaction(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.AccountTransaction{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
