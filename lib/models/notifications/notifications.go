package models

import "time"

// Notification represents a notification message structure
type Notification struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"type:varchar(255)" json:"username"`
	Message   string `gorm:"type:varchar(50)" json:"message"`
	FromRole  string `gorm:"type:varchar(50)" json:"from_role"` // Sender's role
	ToRoles   string `gorm:"type:varchar(50)" json:"to_roles"`  // Recipient roles
	CreatedAt time.Time
}

// ID         uint      `gorm:"primaryKey" json:"id"`
// FirstName  string    `gorm:"type:varchar(50)" validate:"required" json:"first_name"`                  // Max 50 characters
// MiddleName string    `gorm:"type:varchar(50)" json:"middle_name,omitempty"`                           // Optional, max 50 characters
// LastName   string    `gorm:"type:varchar(50)" validate:"required" json:"last_name"`                   // Max 50 characters
// Username   string    `gorm:"type:varchar(255);unique" validate:"required" json:"username"`            // Max 255 characters
// Email      string    `gorm:"type:varchar(255);unique" validate:"required,email" json:"email"`         // Unique and valid email
// Password   string    `gorm:"type:varchar(255)" validate:"required" json:"-"`                          // Max 255 characters
// Role       string    `gorm:"type:varchar(50)" validate:"required,oneof=admin user guest" json:"role"` // Max 50 characters
// VerifiedAt time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
