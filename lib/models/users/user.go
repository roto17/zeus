package models

import "time"

type User struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	FirstName  string     `gorm:"type:varchar(50)" validate:"required" json:"first_name"`                  // Max 50 characters
	MiddleName string     `gorm:"type:varchar(50)" json:"middle_name,omitempty"`                           // Optional, max 50 characters
	LastName   string     `gorm:"type:varchar(50)" validate:"required" json:"last_name"`                   // Max 50 characters
	Username   string     `gorm:"type:varchar(255);unique" validate:"required" json:"username"`            // Max 255 characters
	Email      string     `gorm:"type:varchar(255);unique" validate:"required,email" json:"email"`         // Unique and valid email
	Password   string     `gorm:"type:varchar(255)" validate:"required" json:"-"`                          // Max 255 characters
	Role       string     `gorm:"type:varchar(50)" validate:"required,oneof=admin user guest" json:"role"` // Max 50 characters
	VerifiedAt *time.Time `gorm:"type:timestamp" json:"verified_at,omitempty"`
}

type CreateUserInput struct {
	FirstName  string     `json:"first_name"`            // Max 50 characters
	MiddleName string     `json:"middle_name,omitempty"` // Optional, max 50 characters
	LastName   string     `json:"last_name"`             // Max 50 characters
	Username   string     `json:"username"`              // Max 255 characters
	Email      string     `json:"email"`                 // Unique and valid email
	Password   string     `json:"Password"`              // Max 255 characters
	Role       string     `json:"role"`                  // Max 50 characters
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

type LoginUserInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// func (LoginUserInput) TableName() string {
// 	return "users" // Replace this with your desired table name
// }
