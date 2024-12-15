package models

import (
	"time"

	model_user "github.com/roto17/zeus/lib/models/users"
)

// Notification represents a notification message structure
type ProductCategory struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Description string          `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	UserID      uint            `gorm:"not null;index" validate:"required" json:"user_id"` // Foreign key to the Category table
	User        model_user.User `gorm:"foreignKey:UserID" validate:"-"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductCategoryEncrypted struct {
	ID          string `validate:"required" json:"id"`
	Description string ` json:"description"`
	// UserID      string `validate:"required" json:"user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// func (pg *ProductCategory) GetCompany() model_company.Company {
// 	return pg.User.Company
// }

// // Define a scope for filtering by company_id
// func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Where("company_id = ?", companyID)
// 	}
// }
