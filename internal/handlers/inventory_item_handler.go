package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/response"
	"farm/internal/validator"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewInventoryItemHandler(db *gorm.DB) *InventoryItemHandler {
	return &InventoryItemHandler{repo: repositories.NewInventoryItemRepository(db)}
}

type InventoryItemHandler struct {
	repo *repositories.InventoryItemRepository
}

func (h *InventoryItemHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	categoryID := c.Query("category_id")
	items, total, err := h.repo.List(categoryID, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = items
	return c.JSON(resp)
}

func (h *InventoryItemHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	item, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}
	return c.JSON(item)
}

func (h *InventoryItemHandler) Create(c fiber.Ctx) error {
	var item models.InventoryItem
	if err := validator.Body(c, &item); err != nil {
		return nil
	}
	item.CreatedBy = middleware.GetUserID(c)
	item.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&item); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&item)
	return c.Status(201).JSON(item)
}

func (h *InventoryItemHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	item, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}
	var input models.InventoryItem
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.repo.Update(item, map[string]interface{}{
		"category_id":    input.CategoryID,
		"name":           input.Name,
		"sku":            input.SKU,
		"unit":           input.Unit,
		"purchase_price": input.PurchasePrice,
		"selling_price":  input.SellingPrice,
		"updated_by":     middleware.GetUserID(c),
	}); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(item)
	return c.JSON(item)
}

func (h *InventoryItemHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
