package models

import "time"

type PregnancyStatus string

const (
	PregnancyStatusMated       PregnancyStatus = "Mated"
	PregnancyStatusPregnant    PregnancyStatus = "Pregnant"
	PregnancyStatusDelivered   PregnancyStatus = "Delivered"
	PregnancyStatusAborted     PregnancyStatus = "Aborted"
	PregnancyStatusMiscarriage PregnancyStatus = "Miscarriage"
	PregnancyStatusFailed      PregnancyStatus = "Failed"
)

type AnimalPregnancy struct {
	ID                  uint            `gorm:"primaryKey" json:"id"`
	AnimalID            uint            `gorm:"not null" json:"animal_id" validate:"required"`
	BreederID           *uint           `json:"breeder_id"`
	MatingDate          string          `gorm:"type:date;not null" json:"mating_date" validate:"required"`
	ExpectedDueDate     string          `gorm:"type:date;not null" json:"expected_due_date" validate:"required"`
	ActualBirthDate     *string         `gorm:"type:date" json:"actual_birth_date"`
	Status              PregnancyStatus `gorm:"size:20;not null;default:Mated" json:"status" validate:"required,oneof=Mated Pregnant Delivered Aborted Miscarriage Failed"`
	Note                *string         `gorm:"type:text" json:"note"`
	NumberOfChildren    *int            `json:"number_of_children"`
	NumberOfMaleChildren   *int          `json:"number_of_male_children"`
	NumberOfFemaleChildren *int          `json:"number_of_female_children"`
	NumberOfDeadChildren   *int          `json:"number_of_dead_children"`
	CreatedBy           uint            `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy           uint            `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`

	Animal  Animal  `gorm:"foreignKey:AnimalID" json:"animal"`
	Breeder *Animal `gorm:"foreignKey:BreederID" json:"breeder"`
}