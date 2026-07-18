package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DateString string

func (d *DateString) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*d = DateString(v.Format("2006-01-02"))
		return nil
	case string:
		*d = DateString(v)
		return nil
	case []byte:
		*d = DateString(string(v))
		return nil
	case nil:
		*d = ""
		return nil
	default:
		return fmt.Errorf("cannot scan %T into DateString", value)
	}
}

func (d DateString) Value() (driver.Value, error) {
	return string(d), nil
}

type AnimalWeightHistory struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	AnimalID   uint       `gorm:"not null" json:"animal_id" validate:"required"`
	Weight     float64    `gorm:"type:decimal(8,2);not null" json:"weight" validate:"required,gt=0"`
	RecordDate DateString `gorm:"type:date;not null" json:"record_date" validate:"required"`
	Remarks    *string    `gorm:"type:text" json:"remarks"`
	CreatedBy  uint       `gorm:"not null;default:0" json:"created_by"`
	UpdatedBy  uint       `gorm:"not null;default:0" json:"updated_by"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Animal *Animal `gorm:"foreignKey:AnimalID" json:"-"`
}
