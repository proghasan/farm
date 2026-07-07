package handlers

import (
	"farm/database"
	"farm/middleware"
	"math"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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

func paginate(c *fiber.Ctx, tx *gorm.DB, dest interface{}) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
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

func setCreatedBy(c *fiber.Ctx, tx *gorm.DB) *gorm.DB {
	uid := middleware.GetUserID(c)
	return tx.Set("created_by", uid)
}

func setUpdatedBy(c *fiber.Ctx, tx *gorm.DB) *gorm.DB {
	uid := middleware.GetUserID(c)
	return tx.Set("updated_by", uid)
}

func useDB(c *fiber.Ctx) *gorm.DB {
	return database.DB
}

type handlerFunc func(c *fiber.Ctx, db *gorm.DB) error

func wrapHandler(fn handlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fn(c, useDB(c))
	}
}
