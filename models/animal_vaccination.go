package models

import "time"

type AnimalVaccination struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	AnimalID        uint    `gorm:"not null" json:"animal_id" validate:"required"`
	VaccineID       uint    `gorm:"not null" json:"vaccine_id" validate:"required"`
	VaccinationDate string  `gorm:"type:date;not null" json:"vaccination_date" validate:"required"`
	NextDueDate    *string   `gorm:"type:date" json:"next_due_date,omitempty"`
	DoctorName     *string   `gorm:"size:150" json:"doctor_name,omitempty"`
	Remarks        *string   `gorm:"type:text" json:"remarks,omitempty"`
	CreatedBy      uint      `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy      uint      `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Animal  Animal  `gorm:"foreignKey:AnimalID" json:"animal,omitempty"`
	Vaccine Vaccine `gorm:"foreignKey:VaccineID" json:"vaccine,omitempty"`
}
