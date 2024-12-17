package models

import (
	"time"

	model_category "github.com/roto17/zeus/lib/models/productcategories"
	model_user "github.com/roto17/zeus/lib/models/users"
	"gorm.io/gorm"
)

// Notification represents a notification message structure
type Product struct {
	ID          uint                           `gorm:"primaryKey" json:"id"`
	Description string                         `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	CategoryID  uint                           `gorm:"not null;index" validate:"required" json:"category_id"` // Foreign key to the Category table
	Category    model_category.ProductCategory `gorm:"foreignKey:CategoryID" validate:"-"`                    // Association to the Category
	UserID      uint                           `gorm:"not null;index" validate:"-" json:"user_id"`            // Foreign key to the Category table
	User        model_user.User                `gorm:"foreignKey:UserID" validate:"-"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductEncrypted struct {
	ID          string `validate:"required" json:"id"`
	Description string `validate:"required" json:"description"`
	// QRCode      string `validate:"required" json:"qr_code"`
	CategoryID string `validate:"required" json:"category_id"` // Foreign key to the User table
	UserID     string `validate:"required" json:"user_id"`
	// CreatedAt   time.Time
	// UpdatedAt   time.Time
}

// // Notification represents a notification message structure
// type ProductEncrypted struct {
// 	ID          string                                  `validate:"required" json:"id"`
// 	Description string                                  `validate:"required" json:"description"`
// 	QRCode      string                                  `validate:"required" json:"qr_code"`
// 	CategoryID  string                                  `validate:"required" json:"category_id"` // Foreign key to the User table
// 	Category    model_category.ProductCategoryEncrypted `json:"-"`                               // Association to the User
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

// func (p *Product) GetCompany() model_company.Company {
// 	return p.Company
// }

// Define a scope for filtering by company_id
// func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Where("company_id = ?", companyID)
// 	}
// }

func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Join the 'users' table and filter by 'company_id'
		return db.Joins("JOIN users ON users.id = products.user_id").
			Where("users.company_id = ?", companyID)
	}
}
