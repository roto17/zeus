package models

import (
	"time"

	model_company "github.com/roto17/zeus/lib/models/companies"
	model_category "github.com/roto17/zeus/lib/models/productcategories"
	"gorm.io/gorm"
)

// Notification represents a notification message structure
type Product struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Description string `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	// QRCode      string                         `gorm:"type:varchar(255)" json:"qr_code"`
	CategoryID uint                           `gorm:"not null;index" json:"category_id"` // Foreign key to the Category table
	Category   model_category.ProductCategory `gorm:"foreignKey:CategoryID"`             // Association to the Category

	CompanyID uint                   `gorm:"not null;index" json:"company_id"` // Foreign key to the Category table
	Company   *model_company.Company `gorm:"foreignKey:CompanyID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductEncrypted struct {
	ID          string `validate:"required" json:"id"`
	Description string `validate:"required" json:"description"`
	QRCode      string `validate:"required" json:"qr_code"`
	CategoryID  string `validate:"required" json:"category_id"` // Foreign key to the User table
	CompanyID   string `validate:"required" json:"company_id"`
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

func (p *Product) GetCompany() *model_company.Company {
	return p.Company
}

// Define a scope for filtering by company_id
func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("company_id = ?", companyID)
	}
}
