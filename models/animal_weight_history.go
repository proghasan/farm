package models

import "time"

type AnimalWeightHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AnimalID   uint      `gorm:"not null" json:"animal_id" validate:"required"`
	Weight     float64   `gorm:"type:decimal(8,2);not null" json:"weight" validate:"required,gt=0"`
	RecordDate string    `gorm:"type:date;not null" json:"record_date" validate:"required"`
	Remarks    *string   `gorm:"type:text" json:"remarks,omitempty"`
	CreatedBy  uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy  uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Animal Animal `gorm:"foreignKey:AnimalID" json:"animal,omitempty"`
}
