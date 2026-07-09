package models

import (
	"time"

	"gorm.io/gorm"
)

type Animal struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TagNo         string         `gorm:"size:50;not null;unique" json:"tag_no" validate:"required,min=1,max=50"`
	SpeciesID     uint           `gorm:"not null" json:"species_id" validate:"required"`
	BreedID       *uint          `json:"breed_id"`
	FatherID      *uint          `json:"father_id"`
	MotherID      *uint          `json:"mother_id"`
	Gender        string         `gorm:"size:10;not null" json:"gender" validate:"required,oneof=Male Female"`
	BirthDate     *string        `json:"birth_date"`
	PurchaseDate  *string        `json:"purchase_date"`
	PurchasePrice float64        `gorm:"type:decimal(12,2);default:0" json:"purchase_price"`
	CurrentWeight *float64       `gorm:"type:decimal(8,2)" json:"current_weight"`
	Color         *string        `gorm:"size:100" json:"color"`
	Status        string         `gorm:"size:20;default:Healthy" json:"status"`
	Remarks       *string        `gorm:"type:text" json:"remarks"`
	CreatedBy     uint           `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy     uint           `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Species           Species              `gorm:"foreignKey:SpeciesID" json:"species"`
	Breed             *Breed               `gorm:"foreignKey:BreedID" json:"breed"`
	Father            *Animal              `gorm:"foreignKey:FatherID" json:"father"`
	Mother            *Animal              `gorm:"foreignKey:MotherID" json:"mother"`
	WeightHistories   []AnimalWeightHistory `gorm:"foreignKey:AnimalID" json:"weight_histories"`
	AnimalVaccinations []AnimalVaccination  `gorm:"foreignKey:AnimalID" json:"vaccinations"`
}
