package models

import (
	"time"

	"gorm.io/gorm"
)

type Animal struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TagNo         string         `gorm:"size:50;not null;unique" json:"tag_no" validate:"required,min=1,max=50"`
	Name          *string        `gorm:"size:150" json:"name,omitempty"`
	SpeciesID     uint           `gorm:"not null" json:"species_id" validate:"required"`
	BreedID       *uint          `json:"breed_id,omitempty"`
	Gender        string         `gorm:"size:10;not null" json:"gender" validate:"required,oneof=Male Female"`
	BirthDate     *string        `json:"birth_date,omitempty"`
	PurchaseDate  *string        `json:"purchase_date,omitempty"`
	PurchasePrice float64        `gorm:"type:decimal(12,2);default:0" json:"purchase_price"`
	CurrentWeight *float64       `gorm:"type:decimal(8,2)" json:"current_weight,omitempty"`
	Color         *string        `gorm:"size:100" json:"color,omitempty"`
	Status        string         `gorm:"size:20;default:Healthy" json:"status"`
	Remarks       *string        `gorm:"type:text" json:"remarks,omitempty"`
	CreatedBy     uint           `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy     uint           `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Species           Species              `gorm:"foreignKey:SpeciesID" json:"species,omitempty"`
	Breed             *Breed               `gorm:"foreignKey:BreedID" json:"breed,omitempty"`
	WeightHistories   []AnimalWeightHistory `gorm:"foreignKey:AnimalID" json:"weight_histories,omitempty"`
	AnimalVaccinations []AnimalVaccination  `gorm:"foreignKey:AnimalID" json:"vaccinations,omitempty"`
}
