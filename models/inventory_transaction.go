package models

import "time"

type InventoryTransaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	InventoryItemID uint      `gorm:"not null" json:"inventory_item_id" validate:"required"`
	TransactionType string    `gorm:"size:20;not null" json:"transaction_type" validate:"required,oneof=Purchase Sale Consumption Adjustment Return Damage"`
	Quantity        float64   `gorm:"type:decimal(12,2);not null" json:"quantity" validate:"required"`
	TransactionDate string    `gorm:"type:date;not null" json:"transaction_date" validate:"required"`
	Remarks         *string   `gorm:"type:text" json:"remarks"`
	CreatedBy       uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy       uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	InventoryItem InventoryItem `gorm:"foreignKey:InventoryItemID" json:"inventory_item"`
}
