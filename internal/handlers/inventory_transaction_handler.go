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

func NewInventoryTransactionHandler(db *gorm.DB) *InventoryTransactionHandler {
	return &InventoryTransactionHandler{repo: repositories.NewInventoryTransactionRepository(db)}
}

type InventoryTransactionHandler struct {
	repo *repositories.InventoryTransactionRepository
}

func (h *InventoryTransactionHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	itemID := c.Query("inventory_item_id")
	txnType := c.Query("transaction_type")
	txns, total, err := h.repo.List(itemID, txnType, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = txns
	return c.JSON(resp)
}

func (h *InventoryTransactionHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	txn, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	return c.JSON(txn)
}

func (h *InventoryTransactionHandler) Create(c fiber.Ctx) error {
	var req request.CreateInventoryTransactionRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	txn := models.InventoryTransaction{
		InventoryItemID: req.InventoryItemID,
		TransactionType: req.TransactionType,
		Quantity:        req.Quantity,
		TransactionDate: req.TransactionDate,
		Remarks:         req.Remarks,
	}
	txn.CreatedBy = middleware.GetUserID(c)
	txn.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&txn); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(&txn)
	return c.Status(201).JSON(txn)
}

func (h *InventoryTransactionHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	txn, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}
	var req request.UpdateInventoryTransactionRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.TransactionType != nil {
		updates["transaction_type"] = *req.TransactionType
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.TransactionDate != nil {
		updates["transaction_date"] = *req.TransactionDate
	}
	if req.Remarks != nil {
		updates["remarks"] = *req.Remarks
	}
	if err := h.repo.Update(txn, updates); err != nil {
		return validator.HandleDBError(c, err)
	}
	h.repo.Preload(txn)
	return c.JSON(txn)
}

func (h *InventoryTransactionHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
