// models/token.go

package models

import (
	"time"

	model_user "github.com/roto17/zeus/lib/models/users"
)

type Token struct {
	ID         uint            `gorm:"primaryKey"`
	Token      string          `gorm:"type:text"`
	UserID     uint            `gorm:"not null;index"`    // Foreign key to the User table
	User       model_user.User `gorm:"foreignKey:UserID"` // Association to the User
	IPAddress  string          `gorm:"type:varchar(45)"`  // To accommodate both IPv4 and IPv6 addresses
	DeviceName string          `gorm:"type:varchar(255)"` // Name of the device
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
