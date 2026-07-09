package handlers

import (
	"farm/middleware"
	"farm/models"
	speciesreq "farm/requests/species"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ListSpecies(c fiber.Ctx, db *gorm.DB) error {
	var species []models.Species
	tx := db.Model(&models.Species{}).Preload("User").Order("id desc")
	if s := c.Query("search"); s != "" {
		tx = tx.Where("name LIKE ?", "%"+s+"%")
	}
	return paginate(c, tx, &species)
}

func GetSpecies(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var species models.Species
	if err := db.Preload("Breeds").First(&species, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
	}
	return c.JSON(species)
}

func CreateSpecies(c fiber.Ctx, db *gorm.DB) error {
	var req speciesreq.CreateSpeciesRequest
	if err := req.FromContext(c); err != nil {
		return nil
	}
	species := models.Species{
		Name: req.Name,
	}
	species.CreatedBy = middleware.GetUserID(c)
	species.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&species).Error; err != nil {
		return handleError(c, err)
	}
	return c.Status(201).JSON(species)
}

func UpdateSpecies(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	var species models.Species
	if err := db.First(&species, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
	}
	var req speciesreq.UpdateSpeciesRequest
	if err := req.FromContext(c); err != nil {
		return nil
	}
	species.Name = req.Name
	species.UpdatedBy = middleware.GetUserID(c)
	if err := db.Save(&species).Error; err != nil {
		return handleError(c, err)
	}
	return c.JSON(species)
}

func DeleteSpecies(c fiber.Ctx, db *gorm.DB) error {
	id := fiber.Params[int](c, "id", 0)
	if err := db.Delete(&models.Species{}, id).Error; err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(204)
}
