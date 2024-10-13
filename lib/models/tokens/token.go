// models/token.go

package models

import (
	"time"

	model_user "github.com/roto17/zeus/lib/models/users"
)

type Token struct {
	ID     uint            `gorm:"primaryKey"`
	Token  string          `gorm:"type:text"`
	UserID uint            `gorm:"not null"`          // Foreign key to the User table
	User   model_user.User `gorm:"foreignKey:UserID"` // Association to the User
	// ExpiresAt time.Time `gorm:"index"` // Index for faster queries
	CreatedAt time.Time
	UpdatedAt time.Time
}
