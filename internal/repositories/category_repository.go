package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) List(page, perPage int) ([]models.InventoryCategory, int64, error) {
	var categories []models.InventoryCategory
	tx := r.DB.Model(&models.InventoryCategory{}).Preload("Items")
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&categories).Error
	return categories, total, err
}

func (r *CategoryRepository) GetByID(id uint) (*models.InventoryCategory, error) {
	var cat models.InventoryCategory
	err := r.DB.Preload("Items").First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) Create(cat *models.InventoryCategory) error {
	return r.DB.Create(cat).Error
}

func (r *CategoryRepository) Update(cat *models.InventoryCategory, updates map[string]interface{}) error {
	return r.DB.Model(cat).Updates(updates).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.DB.Delete(&models.InventoryCategory{}, id).Error
}
