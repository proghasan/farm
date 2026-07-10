package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type InventoryTransactionRepository struct {
	DB *gorm.DB
}

func NewInventoryTransactionRepository(db *gorm.DB) *InventoryTransactionRepository {
	return &InventoryTransactionRepository{DB: db}
}

func (r *InventoryTransactionRepository) List(itemID, txnType string, page, perPage int) ([]models.InventoryTransaction, int64, error) {
	var txns []models.InventoryTransaction
	tx := r.DB.Model(&models.InventoryTransaction{}).Preload("InventoryItem.Category")
	if itemID != "" {
		tx = tx.Where("inventory_item_id = ?", itemID)
	}
	if txnType != "" {
		tx = tx.Where("transaction_type = ?", txnType)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Order("transaction_date DESC").Find(&txns).Error
	return txns, total, err
}

func (r *InventoryTransactionRepository) GetByID(id uint) (*models.InventoryTransaction, error) {
	var txn models.InventoryTransaction
	err := r.DB.Preload("InventoryItem.Category").First(&txn, id).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *InventoryTransactionRepository) Create(txn *models.InventoryTransaction) error {
	return r.DB.Create(txn).Error
}

func (r *InventoryTransactionRepository) Update(txn *models.InventoryTransaction, updates map[string]interface{}) error {
	return r.DB.Model(txn).Updates(updates).Error
}

func (r *InventoryTransactionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.InventoryTransaction{}, id).Error
}

func (r *InventoryTransactionRepository) Preload(txn *models.InventoryTransaction) error {
	return r.DB.Preload("InventoryItem.Category").First(txn, txn.ID).Error
}
