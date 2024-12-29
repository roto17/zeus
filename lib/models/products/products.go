package models

import (
	"time"

	model_category "github.com/roto17/zeus/lib/models/productcategories"
	model_user "github.com/roto17/zeus/lib/models/users"
	"gorm.io/gorm"
)

// Product structure
type Product struct {
	ID          uint   `gorm:"primaryKey" json:"id,omitempty"`
	Description string `gorm:"not null;type:varchar(50);unique" validate:"required" json:"description,omitempty"`
	CategoryID  uint   `gorm:"not null;index" json:"category_id,omitempty"` // Foreign key to the Category table
	// Category    model_category.ProductCategory `gorm:"foreignKey:CategoryID" validate:"-" json:"category"`    // Association to the Category
	// Category  model_category.ProductCategoryResponse `gorm:"foreignKey:CategoryID" validate:"-" json:"category,omitempty"` // Association to the Category
	UserID uint `gorm:"not null;index" validate:"required" json:"user_id,omitempty"` // Foreign key to the Category table
	// User      model_user.UserResponse                `gorm:"foreignKey:UserID" validate:"-" json:"user,omitempty"`
	Unit         string    `gorm:"type:varchar(50)" validate:"required,oneof=Piece Metre KG" json:"unit,omitempty"`
	Quantity     int64     `gorm:"type:int;default:0" json:"quantity,omitempty"`
	BuyingPrice  float64   `gorm:"type:decimal(10,2)" validate:"required" json:"buying_price,omitempty"`  // Buying price
	SellingPrice float64   `gorm:"type:decimal(10,2)" validate:"required" json:"selling_price,omitempty"` // Selling price
	Weight       float64   `gorm:"type:decimal(5,2)" json:"selling_weight"`                               // Weight in percentage (e.g., 10.00 represents 10%)
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type ProductEncrypted struct {
	ID          string `validate:"required" json:"id,omitempty"`
	Description string `validate:"required" json:"description,omitempty"`
	// QRCode      string `validate:"required" json:"qr_code"`
	CategoryID   string  `validate:"required" json:"category_id,omitempty"` // Foreign key to the User table
	UserID       string  `validate:"required" json:"user_id,omitempty"`
	Unit         string  `validate:"required,oneof=Piece Metre KG" json:"unit,omitempty"`
	Quantity     int64   `json:"quantity,omitempty"`
	BuyingPrice  float64 `validate:"required" json:"buying_price,omitempty"`  // Buying price
	SellingPrice float64 `validate:"required" json:"selling_price,omitempty"` // Selling price
	Weight       float64 `json:"selling_weight,omitempty"`                    // Weight in percentage (e.g., 10.00 represents 10%)
	// CreatedAt   time.Time
	// UpdatedAt   time.Time
}

// Product structure
type ProductResponse struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Description string `gorm:"type:varchar(50);unique" validate:"-" json:"description"`
	CategoryID  uint   `gorm:"not null;index" validate:"-" json:"-"` // Foreign key to the Category table
	// Category    model_category.ProductCategory `gorm:"foreignKey:CategoryID" validate:"-" json:"category"`    // Association to the Category
	Category     model_category.ProductCategoryResponse `gorm:"foreignKey:CategoryID" validate:"-" json:"category"` // Association to the Category
	UserID       uint                                   `gorm:"not null;index" validate:"-" json:"-"`               // Foreign key to the Category table
	User         model_user.UserResponse                `gorm:"foreignKey:UserID" validate:"-" json:"user"`
	Unit         string                                 `gorm:"type:varchar(50)" validate:"-" json:"unit"`
	Quantity     int64                                  `gorm:"type:int;default:0" json:"quantity"`
	BuyingPrice  float64                                `gorm:"type:decimal(10,2)" validate:"-" json:"buying_price"`  // Buying price
	SellingPrice float64                                `gorm:"type:decimal(10,2)" validate:"-" json:"selling_price"` // Selling price
	Weight       float64                                `gorm:"type:decimal(5,2)" json:"selling_weight"`              // Weight in percentage (e.g., 10.00 represents 10%)
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Product structure
type ProductPatch struct {
	ID          uint   `gorm:"primaryKey" json:"id,omitempty"`
	Description string `gorm:"not null;type:varchar(50);unique" validate:"-" json:"description,omitempty"`
	CategoryID  uint   `gorm:"not null;index" json:"category_id,omitempty"` // Foreign key to the Category table
	// Category    model_category.ProductCategory `gorm:"foreignKey:CategoryID" validate:"-" json:"category"`    // Association to the Category
	// Category  model_category.ProductCategoryResponse `gorm:"foreignKey:CategoryID" validate:"-" json:"category,omitempty"` // Association to the Category
	UserID uint `gorm:"not null;index" validate:"-" json:"user_id,omitempty"` // Foreign key to the Category table
	// User      model_user.UserResponse                `gorm:"foreignKey:UserID" validate:"-" json:"user,omitempty"`
	Unit         string    `gorm:"type:varchar(50)" validate:"oneof=Piece Metre KG" json:"unit,omitempty"`
	Quantity     int64     `gorm:"type:int;default:0" json:"quantity,omitempty"`
	BuyingPrice  float64   `gorm:"type:decimal(10,2)" validate:"-" json:"buying_price,omitempty"`  // Buying price
	SellingPrice float64   `gorm:"type:decimal(10,2)" validate:"-" json:"selling_price,omitempty"` // Selling price
	Weight       float64   `gorm:"type:decimal(5,2)" validate:"-" json:"selling_weight"`           // Weight in percentage (e.g., 10.00 represents 10%)
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

func (ProductResponse) TableName() string {
	return "products" // Replace this with your desired table name
}

func (ProductPatch) TableName() string {
	return "products" // Replace this with your desired table name
}

// func (p *Product) GetCompany() model_company.Company {
// 	return p.Company
// }

func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Join the 'users' table and filter by 'company_id'
		return db.Joins("JOIN users ON users.id = products.user_id").
			Where("users.company_id = ?", companyID)
	}
}
