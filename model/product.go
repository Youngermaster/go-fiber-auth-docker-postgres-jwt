package model

import "gorm.io/gorm"

// Product struct
type Product struct {
	gorm.Model
	Title       string `gorm:"not null;size:255" json:"title" validate:"required,min=1,max=255"`
	Description string `gorm:"not null;type:text" json:"description" validate:"required,min=1,max=2000"`
	Amount      int    `gorm:"not null" json:"amount" validate:"required,min=0"`
	UserID      uint   `gorm:"not null;index" json:"user_id"`
}
