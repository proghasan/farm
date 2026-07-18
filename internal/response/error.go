package response

import (
	"errors"
	"strings"

	"farm/pkg/validator"

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
			"errors":  []string{"Resource not found"},
		})

	case errors.As(err, new(*fiber.Error)):
		var e *fiber.Error
		errors.As(err, &e)

		return c.Status(e.Code).JSON(fiber.Map{
			"success": false,
			"errors":  []string{e.Message},
		})

	case errors.As(err, new(*validator.ValidationError)):
		var ve *validator.ValidationError
		errors.As(err, &ve)

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  ve.Errors,
		})

	case isForeignKeyError(err):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"errors":  []string{"Cannot delete: resource is in use."},
		})

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  []string{"Internal server error"},
		})
	}
}

func isForeignKeyError(err error) bool {
	return strings.Contains(err.Error(), "foreign key constraint")
}