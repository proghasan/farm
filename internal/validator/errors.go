package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func HandleBindError(c fiber.Ctx, err error) {
	var vErr validator.ValidationErrors
	if errors.As(err, &vErr) {
		WriteErrors(c, vErr)
		return
	}
	c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
}

func WriteErrors(c fiber.Ctx, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Field() + " is " + e.Tag()
		}
		c.Status(422).JSON(fiber.Map{"errors": msgs})
		return
	}
	c.Status(422).JSON(fiber.Map{"errors": []string{err.Error()}})
}

func HandleDBError(c fiber.Ctx, err error) error {
	msg := err.Error()
	if strings.Contains(msg, "Error 1452") {
		parts := strings.Split(msg, "CONSTRAINT `")
		if len(parts) > 1 {
			fk := strings.Split(parts[1], "`")[0]
			return c.Status(422).JSON(fiber.Map{"error": "Referenced record not found (" + fk + ")"})
		}
		return c.Status(422).JSON(fiber.Map{"error": "Referenced record not found"})
	}
	if strings.Contains(msg, "Duplicate entry") {
		return c.Status(409).JSON(fiber.Map{"error": "Duplicate value, record already exists"})
	}
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}
	return c.Status(400).JSON(fiber.Map{"error": msg})
}
