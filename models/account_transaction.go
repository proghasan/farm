package models

import "time"

type AccountTransaction struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AccountHeadID uint      `gorm:"not null" json:"account_head_id" validate:"required"`
	TransactionDate string  `gorm:"type:date;not null" json:"transaction_date" validate:"required"`
	Amount        float64   `gorm:"type:decimal(12,2);not null" json:"amount" validate:"required"`
	PaymentMethod string    `gorm:"size:20;default:Cash" json:"payment_method" validate:"omitempty,oneof=Cash Bank 'Mobile Banking' Other"`
	ReferenceNo   *string   `gorm:"size:100" json:"reference_no,omitempty"`
	Description   *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedBy     uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy     uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	AccountHead AccountHead `gorm:"foreignKey:AccountHeadID" json:"account_head,omitempty"`
}
