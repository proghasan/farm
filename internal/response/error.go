package response

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Resource not found",
		})

	case errors.As(err, new(*fiber.Error)):
		var e *fiber.Error
		errors.As(err, &e)

		return c.Status(e.Code).JSON(fiber.Map{
			"success": false,
			"message": e.Message,
		})

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}
}