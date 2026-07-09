package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var Validate = validator.New()

func WriteValidationErrors(c fiber.Ctx, err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Field() + " is " + e.Tag()
		}
		c.Status(422).JSON(fiber.Map{"errors": msgs})
		return err
	}
	c.Status(422).JSON(fiber.Map{"errors": []string{err.Error()}})
	return err
}
