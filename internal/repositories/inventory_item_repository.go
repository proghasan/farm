package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type InventoryItemRepository struct {
	DB *gorm.DB
}

func NewInventoryItemRepository(db *gorm.DB) *InventoryItemRepository {
	return &InventoryItemRepository{DB: db}
}

func (r *InventoryItemRepository) List(categoryID string, page, perPage int) ([]models.InventoryItem, int64, error) {
	var items []models.InventoryItem
	tx := r.DB.Model(&models.InventoryItem{}).Preload("Category")
	if categoryID != "" {
		tx = tx.Where("category_id = ?", categoryID)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&items).Error
	return items, total, err
}

func (r *InventoryItemRepository) GetByID(id uint) (*models.InventoryItem, error) {
	var item models.InventoryItem
	err := r.DB.Preload("Category").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *InventoryItemRepository) Create(item *models.InventoryItem) error {
	return r.DB.Create(item).Error
}

func (r *InventoryItemRepository) Update(item *models.InventoryItem, updates map[string]interface{}) error {
	return r.DB.Model(item).Updates(updates).Error
}

func (r *InventoryItemRepository) Delete(id uint) error {
	return r.DB.Delete(&models.InventoryItem{}, id).Error
}

func (r *InventoryItemRepository) Preload(item *models.InventoryItem) error {
	return r.DB.Preload("Category").First(item, item.ID).Error
}
