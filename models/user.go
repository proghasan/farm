package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:150;not null" json:"name"`
	Email         *string        `gorm:"size:255;unique" json:"email,omitempty"`
	Phone         *string        `gorm:"size:20;unique" json:"phone,omitempty"`
	Username      *string        `gorm:"size:100;unique" json:"username,omitempty"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	Avatar        *string        `gorm:"size:255" json:"avatar,omitempty"`
	Role          string         `gorm:"size:20;default:Worker" json:"role"`
	Status        string         `gorm:"size:20;default:Active" json:"status"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
	RememberToken *string        `gorm:"size:100" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
