package models

import (
	"time"

	model_category "github.com/roto17/zeus/lib/models/productcategories"
)

// Notification represents a notification message structure
type Product struct {
	ID          uint                           `gorm:"primaryKey" json:"id"`
	Description string                         `gorm:"type:varchar(50);unique" validate:"required" json:"description"`
	QRCode      string                         `gorm:"type:varchar(255)" json:"qr_code"`
	CategoryID  uint                           `gorm:"not null;index" json:"category_id"` // Foreign key to the User table
	Category    model_category.ProductCategory `gorm:"foreignKey:CategoryID"`             // Association to the User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductInput struct {
	// ID          uint   `json:"id"`
	Description string `json:"description"`
	QRCode      string `json:"qr_code"`
	CategoryID  uint   `json:"category_id"` // Foreign key to the User table
	// CreatedAt   time.Time
	// UpdatedAt   time.Time
}

// Notification represents a notification message structure
type ProductEncrypted struct {
	ID          string                                  `validate:"required" json:"id"`
	Description string                                  `validate:"required" json:"description"`
	QRCode      string                                  `validate:"required" json:"qr_code"`
	CategoryID  string                                  `validate:"required" json:"category_id"` // Foreign key to the User table
	Category    model_category.ProductCategoryEncrypted `json:"-"`                               // Association to the User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
