package models

import "time"

type AccountHead struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Type        string    `gorm:"size:20;not null" json:"type" validate:"required,oneof=Income Expense"`
	Name        string    `gorm:"size:150;not null" json:"name" validate:"required,min=1,max=150"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedBy   uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy   uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Transactions []AccountTransaction `gorm:"foreignKey:AccountHeadID" json:"transactions"`
}
