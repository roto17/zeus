package models

import "time"

// Notification represents a notification message structure
type Company struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Description string `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CompanyEncrypted struct {
	ID          string `validate:"required" json:"id"`
	Description string `validate:"required" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
