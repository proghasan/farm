package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type AccountHeadRepository struct {
	DB *gorm.DB
}

func NewAccountHeadRepository(db *gorm.DB) *AccountHeadRepository {
	return &AccountHeadRepository{DB: db}
}

func (r *AccountHeadRepository) List(acctType string, page, perPage int) ([]models.AccountHead, int64, error) {
	var heads []models.AccountHead
	tx := r.DB.Model(&models.AccountHead{})
	if acctType != "" {
		tx = tx.Where("type = ?", acctType)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&heads).Error
	return heads, total, err
}

func (r *AccountHeadRepository) GetByID(id uint) (*models.AccountHead, error) {
	var head models.AccountHead
	err := r.DB.First(&head, id).Error
	if err != nil {
		return nil, err
	}
	return &head, nil
}

func (r *AccountHeadRepository) Create(head *models.AccountHead) error {
	return r.DB.Create(head).Error
}

func (r *AccountHeadRepository) Update(head *models.AccountHead, updates map[string]interface{}) error {
	return r.DB.Model(head).Updates(updates).Error
}

func (r *AccountHeadRepository) Delete(id uint) error {
	return r.DB.Delete(&models.AccountHead{}, id).Error
}
