package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type AccountTransactionRepository struct {
	DB *gorm.DB
}

func NewAccountTransactionRepository(db *gorm.DB) *AccountTransactionRepository {
	return &AccountTransactionRepository{DB: db}
}

func (r *AccountTransactionRepository) List(headID, acctType string, page, perPage int) ([]models.AccountTransaction, int64, error) {
	var txns []models.AccountTransaction
	tx := r.DB.Model(&models.AccountTransaction{}).Preload("AccountHead")
	if headID != "" {
		tx = tx.Where("account_head_id = ?", headID)
	}
	if acctType != "" {
		tx = tx.Joins("AccountHead").Where("AccountHead.type = ?", acctType)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Order("transaction_date DESC").Find(&txns).Error
	return txns, total, err
}

func (r *AccountTransactionRepository) GetByID(id uint) (*models.AccountTransaction, error) {
	var txn models.AccountTransaction
	err := r.DB.Preload("AccountHead").First(&txn, id).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *AccountTransactionRepository) Create(txn *models.AccountTransaction) error {
	return r.DB.Create(txn).Error
}

func (r *AccountTransactionRepository) Update(txn *models.AccountTransaction, updates map[string]interface{}) error {
	return r.DB.Model(txn).Updates(updates).Error
}

func (r *AccountTransactionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.AccountTransaction{}, id).Error
}

func (r *AccountTransactionRepository) Preload(txn *models.AccountTransaction) error {
	return r.DB.Preload("AccountHead").First(txn, txn.ID).Error
}
