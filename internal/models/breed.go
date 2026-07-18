package models

import "time"

type Breed struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SpeciesID uint      `gorm:"not null" json:"species_id" validate:"required"`
	Name      string    `gorm:"size:150;not null" json:"name" validate:"required,min=1,max=150"`
	CreatedBy uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Species Species `gorm:"foreignKey:SpeciesID" json:"species"`
	User    User    `gorm:"foreignKey:CreatedBy" json:"user"`
}
