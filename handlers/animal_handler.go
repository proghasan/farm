package handlers

import (
	"farm/middleware"
	"farm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListAnimals(c *fiber.Ctx, db *gorm.DB) error {
	var animals []models.Animal
	tx := db.Preload("Species").Preload("Breed")
	if speciesID := c.Query("species_id"); speciesID != "" {
		tx = tx.Where("species_id = ?", speciesID)
	}
	if status := c.Query("status"); status != "" {
		tx = tx.Where("status = ?", status)
	}
	if gender := c.Query("gender"); gender != "" {
		tx = tx.Where("gender = ?", gender)
	}
	return paginate(c, tx, &animals)
}

func GetAnimal(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var animal models.Animal
	if err := db.
		Preload("Species").
		Preload("Breed").
		Preload("WeightHistories").
		Preload("AnimalVaccinations.Vaccine").
		First(&animal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
	}
	return c.JSON(animal)
}

func CreateAnimal(c *fiber.Ctx, db *gorm.DB) error {
	var animal models.Animal
	if err := c.BodyParser(&animal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	animal.CreatedBy = middleware.GetUserID(c)
	animal.UpdatedBy = middleware.GetUserID(c)
	if err := db.Create(&animal).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(animal)
}

func UpdateAnimal(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	var animal models.Animal
	if err := db.First(&animal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
	}
	var input models.Animal
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.ID = animal.ID
	input.CreatedBy = animal.CreatedBy
	input.UpdatedBy = middleware.GetUserID(c)
	db.Model(&animal).Updates(input)
	return c.JSON(animal)
}

func DeleteAnimal(c *fiber.Ctx, db *gorm.DB) error {
	id, _ := c.ParamsInt("id")
	if err := db.Delete(&models.Animal{}, id).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
