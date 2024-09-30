// models/token.go

package models

import (
	"time"
)

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:text"`
	ExpiresAt time.Time `gorm:"index"` // Index for faster queries
	CreatedAt time.Time
	UpdatedAt time.Time
}
