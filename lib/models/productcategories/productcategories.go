package models

import (
	"time"

	// model_company "github.com/roto17/zeus/lib/models/companies"
	model_user "github.com/roto17/zeus/lib/models/users"
	"gorm.io/gorm"
)

// Notification represents a notification message structure
type ProductCategory struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Description string          `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	UserID      uint            `gorm:"not null;index" validate:"required" json:"user_id"` // Foreign key to the Category table
	User        model_user.User `gorm:"foreignKey:UserID" validate:"-" json:"user"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductCategoryEncrypted struct {
	ID          string `validate:"required" json:"id"`
	Description string `validate:"required" json:"description"`
	UserID      string ` json:"user_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// func (pg *ProductCategory) GetCompany() model_company.Company {
// 	return pg.User.Company
// }

// FilterByCompanyID filters ProductCategories by the company_id of the User
func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Join the 'users' table and filter by 'company_id'
		return db.Joins("JOIN users ON users.id = product_categories.user_id").
			Where("users.company_id = ?", companyID)
	}
}

type ProductCategoryResponse struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Description string          `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	UserID      uint            `gorm:"not null;index" validate:"required" json:"-"` // Foreign key to the Category table
	User        model_user.User `gorm:"foreignKey:UserID" validate:"-" json:"-"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
}

func (ProductCategoryResponse) TableName() string {
	return "product_categories" // Replace this with your desired table name
}
