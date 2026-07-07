package models

import "time"

type InventoryCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:150;not null;unique" json:"name" validate:"required,min=1,max=150"`
	CreatedBy uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Items []InventoryItem `gorm:"foreignKey:CategoryID" json:"items,omitempty"`
}
