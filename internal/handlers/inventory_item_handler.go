package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
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
	var req request.CreateInventoryItemRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	item := models.InventoryItem{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		SKU:           req.SKU,
		Unit:          req.Unit,
		PurchasePrice: req.PurchasePrice,
		SellingPrice:  req.SellingPrice,
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
	var req request.UpdateInventoryItemRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.Unit != nil {
		updates["unit"] = *req.Unit
	}
	if req.PurchasePrice != nil {
		updates["purchase_price"] = *req.PurchasePrice
	}
	if req.SellingPrice != nil {
		updates["selling_price"] = *req.SellingPrice
	}
	if err := h.repo.Update(item, updates); err != nil {
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
