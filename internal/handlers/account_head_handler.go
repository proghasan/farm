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

func NewAccountHeadHandler(db *gorm.DB) *AccountHeadHandler {
	return &AccountHeadHandler{repo: repositories.NewAccountHeadRepository(db)}
}

type AccountHeadHandler struct {
	repo *repositories.AccountHeadRepository
}

func (h *AccountHeadHandler) List(c fiber.Ctx) error {
	page := fiber.Query[int](c, "page", 1)
	perPage := fiber.Query[int](c, "per_page", 20)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	acctType := c.Query("type")
	heads, total, err := h.repo.List(acctType, page, perPage)
	if err != nil {
		return err
	}
	resp := response.Paginate(page, perPage, total)
	resp.Data = heads
	return c.JSON(resp)
}

func (h *AccountHeadHandler) Get(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	head, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account head not found"})
	}
	return c.JSON(head)
}

func (h *AccountHeadHandler) Create(c fiber.Ctx) error {
	var req request.CreateAccountHeadRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	head := models.AccountHead{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
	}
	head.CreatedBy = middleware.GetUserID(c)
	head.UpdatedBy = middleware.GetUserID(c)
	if err := h.repo.Create(&head); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.Status(201).JSON(head)
}

func (h *AccountHeadHandler) Update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	head, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account head not found"})
	}
	var req request.UpdateAccountHeadRequest
	if err := c.Bind().Body(&req); err != nil {
		validator.HandleBindError(c, err)
		return nil
	}
	updates := map[string]interface{}{
		"updated_by": middleware.GetUserID(c),
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if err := h.repo.Update(head, updates); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.JSON(head)
}

func (h *AccountHeadHandler) Delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id", 0)
	if err := h.repo.Delete(uint(id)); err != nil {
		return validator.HandleDBError(c, err)
	}
	return c.SendStatus(204)
}
