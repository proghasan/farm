package species

import (
	"strings"

	"farm/requests"

	"github.com/gofiber/fiber/v3"
)

type UpdateSpeciesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (r *UpdateSpeciesRequest) FromContext(c fiber.Ctx) error {
	if err := c.Bind().Body(r); err != nil {
		c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		return err
	}
	r.Name = strings.TrimSpace(r.Name)
	if err := requests.Validate.Struct(r); err != nil {
		return requests.WriteValidationErrors(c, err)
	}
	return nil
}
