package handlers

import (
	"farm/internal/middleware"
	"farm/internal/models"
	"farm/internal/repositories"
	"farm/internal/request"
	"farm/internal/response"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewAccountTransactionHandler(db *gorm.DB) *AccountTransactionHandler {
	return &AccountTransactionHandler{repo: repositories.NewAccountTransactionRepository(db)}
}

type AccountTransactionHandler struct {
	repo *repositories.AccountTransactionRepository
}

func (h *AccountTransactionHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	headID := c.Query("account_head_id")
	acctType := c.Query("type")
	txns, total, err := h.repo.List(headID, acctType, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = txns
	return c.JSON(resp)
}

func (h *AccountTransactionHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	txn, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	return c.JSON(txn)
}

func (h *AccountTransactionHandler) Create(c fiber.Ctx) error {
	var req request.CreateAccountTransactionRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	txn := models.AccountTransaction{
		AccountHeadID:   req.AccountHeadID,
		TransactionDate: req.TransactionDate,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		ReferenceNo:     req.ReferenceNo,
		Description:     req.Description,
	}
	txn.CreatedBy = middleware.GetUserID(c)
	txn.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&txn); err != nil {
		return err
	}
	h.repo.Preload(&txn)
	return c.Status(201).JSON(txn)
}

func (h *AccountTransactionHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	txn, err := h.repo.GetByID(uint(id))
	if err != nil {
		return err
	}
	var req request.UpdateAccountTransactionRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.TransactionDate != nil {
		updates["transaction_date"] = *req.TransactionDate
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.PaymentMethod != nil {
		updates["payment_method"] = *req.PaymentMethod
	}
	if req.ReferenceNo != nil {
		updates["reference_no"] = *req.ReferenceNo
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if err := h.repo.Update(txn, updates); err != nil {
		return err
	}
	h.repo.Preload(txn)
	return c.JSON(txn)
}

func (h *AccountTransactionHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return err
	}
	return c.SendStatus(204)
}
