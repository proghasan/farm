package handlers

import (
	"farm/database"
	"farm/middleware"
	"math"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

var validate = validator.New()

type PaginatedRequest struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Search  string `query:"search"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

func paginate(c fiber.Ctx, tx *gorm.DB, dest interface{}) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	var total int64
	tx.Count(&total)

	offset := (page - 1) * perPage
	result := tx.Offset(offset).Limit(perPage).Find(dest)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(PaginatedResponse{
		Data:       dest,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
	})
}

func setCreatedBy(c fiber.Ctx, tx *gorm.DB) *gorm.DB {
	uid := middleware.GetUserID(c)
	return tx.Set("created_by", uid)
}

func setUpdatedBy(c fiber.Ctx, tx *gorm.DB) *gorm.DB {
	uid := middleware.GetUserID(c)
	return tx.Set("updated_by", uid)
}

func useDB(c fiber.Ctx) *gorm.DB {
	return database.DB
}

type handlerFunc func(c fiber.Ctx, db *gorm.DB) error

func wrapHandler(fn handlerFunc) fiber.Handler {
	return func(c fiber.Ctx) error {
		return fn(c, useDB(c))
	}
}

func recordExists(db *gorm.DB, model interface{}, id uint) (bool, error) {
	var count int64
	err := db.Model(model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func handleError(c fiber.Ctx, err error) error {
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
	return c.Status(400).JSON(fiber.Map{"error": msg})
}

func validateBody(c fiber.Ctx, out interface{}) error {
	if err := c.Bind().Body(out); err != nil {
		c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		return err
	}
	if err := validate.Struct(out); err != nil {
		errs := err.(validator.ValidationErrors)
		msgs := make([]string, 0)
		for _, e := range errs {
			msgs = append(msgs, e.Field()+" is "+e.Tag())
		}
		c.Status(422).JSON(fiber.Map{"errors": msgs})
		return err
	}
	return nil
}
