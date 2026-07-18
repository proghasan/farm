package models

import "time"

type Species struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	CreatedBy uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Breeds []Breed `gorm:"foreignKey:SpeciesID" json:"breeds"`
	User *User `gorm:"foreignKey:CreatedBy" json:"user"`
}
