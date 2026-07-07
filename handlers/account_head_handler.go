package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListAccountHeads(c *fiber.Ctx, db *gorm.DB) error {
	var heads []models.AccountHead
	tx := db.Model(&models.AccountHead{})
	if acctType := c.Query("type"); acctType != "" {
		tx = tx.Where("type = ?", acctType)
	}
	return paginate(c, tx, &heads)
}

func GetAccountHead(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var head models.AccountHead
	if err := db.First(&head, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account head not found"})
	}
	return c.JSON(head)
}

func CreateAccountHead(c *fiber.Ctx, db *gorm.DB) error {
	var head models.AccountHead
	if err := validateBody(c, &head); err != nil {
		return err
	}
	head.CreatedBy = middleware.GetUserID(c)
	head.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&head).Error; err != nil {
		return handleError(c, err)
	}
	return c.Status(201).JSON(head)
}

func UpdateAccountHead(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var head models.AccountHead
	if err := db.First(&head, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account head not found"})
	}
	var input models.AccountHead
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = head.ID
	input.CreatedBy = head.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&head).Updates(input)
	return c.JSON(head)
}

func DeleteAccountHead(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.AccountHead{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
