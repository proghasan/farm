package models

import "time"

type Vaccine struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	SpeciesID       uint   `gorm:"not null" json:"species_id"`
	Name            string `gorm:"size:150;not null" json:"name"`
	Description     *string `gorm:"type:text" json:"description,omitempty"`
	Dose            *string `gorm:"size:100" json:"dose,omitempty"`
	MinimumAgeValue uint   `gorm:"not null" json:"minimum_age_value"`
	MinimumAgeUnit  string `gorm:"size:20;not null" json:"minimum_age_unit"`
	IntervalValue   int    `gorm:"not null" json:"interval_value"`
	IntervalUnit    string `gorm:"size:20;not null" json:"interval_unit"`
	IsRepeatable    bool   `gorm:"default:true" json:"is_repeatable"`
	CreatedBy       uint   `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy       uint   `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Species Species `gorm:"foreignKey:SpeciesID" json:"species,omitempty"`
}
