package models

import (
	"time"

	model_category "github.com/roto17/zeus/lib/models/productcategories"
)

// Notification represents a notification message structure
type Product struct {
	ID          uint                           `gorm:"primaryKey" json:"id"`
	Description string                         `gorm:"type:varchar(50)" json:"username"`
	QRCode      string                         `gorm:"type:varchar(255)" json:"qr_code"`
	CategoryID  uint                           `gorm:"not null;index"`        // Foreign key to the User table
	Category    model_category.ProductCategory `gorm:"foreignKey:CategoryID"` // Association to the User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
