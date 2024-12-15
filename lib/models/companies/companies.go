package models

import (
	"time"

	"gorm.io/gorm"
)

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

// Define an interface that models with a Company field must implement
type CompanyProvider interface {
	GetCompany() *Company
	FilterByCompanyID(companyID uint) *gorm.DB
}

func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", companyID)
	}
}
