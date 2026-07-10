package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var Validate = validator.New()

func Body(c fiber.Ctx, out interface{}) error {
	if err := c.Bind().Body(out); err != nil {
		c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		return err
	}
	if err := Validate.Struct(out); err != nil {
		WriteErrors(c, err)
		return err
	}
	return nil
}
