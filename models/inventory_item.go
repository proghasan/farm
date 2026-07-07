package models

import "time"

type InventoryItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CategoryID    uint      `gorm:"not null" json:"category_id" validate:"required"`
	Name          string    `gorm:"size:200;not null" json:"name" validate:"required,min=1,max=200"`
	SKU           *string   `gorm:"size:100" json:"sku,omitempty"`
	Unit          string    `gorm:"size:50;not null" json:"unit" validate:"required,min=1,max=50"`
	PurchasePrice float64   `gorm:"type:decimal(12,2);default:0" json:"purchase_price"`
	SellingPrice  float64   `gorm:"type:decimal(12,2);default:0" json:"selling_price"`
	CreatedBy     uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy     uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Category     InventoryCategory      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Transactions []InventoryTransaction `gorm:"foreignKey:InventoryItemID" json:"transactions,omitempty"`
}
