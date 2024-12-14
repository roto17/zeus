package models

import (
	"time"

	// "gorm.io/gorm"
	model_company "github.com/roto17/zeus/lib/models/companies"
	"gorm.io/gorm"
)

type User struct {
	ID             uint                   `gorm:"primaryKey" json:"id"`
	FirstName      string                 `gorm:"type:varchar(50)" validate:"required" json:"first_name"`                  // Max 50 characters
	MiddleName     string                 `gorm:"type:varchar(50)" json:"middle_name,omitempty"`                           // Optional, max 50 characters
	LastName       string                 `gorm:"type:varchar(50)" validate:"required" json:"last_name"`                   // Max 50 characters
	Username       string                 `gorm:"type:varchar(255);unique" validate:"required" json:"username"`            // Max 255 characters
	Email          string                 `gorm:"type:varchar(255);unique" validate:"required,email" json:"email"`         // Unique and valid email
	Password       string                 `gorm:"type:varchar(255)" validate:"required" json:"-"`                          // Max 255 characters
	Role           string                 `gorm:"type:varchar(50)" validate:"required,oneof=admin user guest" json:"role"` // Max 50 characters
	CompanyID      uint                   `gorm:"not null;index" json:"company_id"`                                        // Foreign key to the Category table
	Company        *model_company.Company `gorm:"foreignKey:CompanyID"`
	VerifiedAt     time.Time              `gorm:"type:timestamp" json:"verified_at,omitempty"`
	VerifiedMethod string                 `gorm:"type:varchar(50)" json:"verified_method,omitempty"`
}

// type CreateUserInput struct {
// 	FirstName  string `json:"first_name"`            // Max 50 characters
// 	MiddleName string `json:"middle_name,omitempty"` // Optional, max 50 characters
// 	LastName   string `json:"last_name"`             // Max 50 characters
// 	Username   string `json:"username"`              // Max 255 characters
// 	Email      string `json:"email"`                 // Unique and valid email
// 	Password   string `json:"Password"`              // Max 255 characters
// 	Role       string `json:"role"`                  // Max 50 characters
// 	CompanyID  string `json:"company_id"`            // Foreign key to the Category table
// 	// VerifiedAt time.Time `json:"verified_at,omitempty"`
// }

type LoginUserInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type EncryptedUser struct {
	ID             string                 `gorm:"-" json:"id"`
	FirstName      string                 `gorm:"-" json:"first_name"`  // Max 50 characters
	MiddleName     string                 `gorm:"-" json:"middle_name"` // Optional, max 50 characters
	LastName       string                 `gorm:"-"  json:"last_name"`  // Max 50 characters
	Username       string                 `gorm:"-"  json:"username"`   // Max 255 characters
	Email          string                 `gorm:"-"  json:"email"`      // Unique and valid email
	Password       string                 `gorm:"-"  json:"Password"`   // Max 255 characters
	Role           string                 `gorm:"-"  json:"role"`       // Max 50 characters
	CompanyID      string                 `gorm:"-" json:"company_id"`  // Foreign key to the Category table
	Company        *model_company.Company `gorm:"-"`
	VerifiedAt     time.Time              `gorm:"-" json:"verified_at,omitempty"`
	VerifiedMethod string                 `gorm:"-" json:"verified_method,omitempty"`
}

type UserUpdateModel struct {
	ID             uint                   `gorm:"primaryKey"`
	FirstName      string                 `gorm:"type:varchar(50)"`  // Max 50 characters
	MiddleName     string                 `gorm:"type:varchar(50)"`  // Optional, max 50 characters
	LastName       string                 `gorm:"type:varchar(50)"`  // Max 50 characters
	Username       string                 `gorm:"type:varchar(255)"` // Max 255 characters
	Email          string                 `gorm:"type:varchar(255)"` // Unique and valid email
	Password       string                 `gorm:"type:varchar(255)"` // Max 255 characters
	Role           string                 `gorm:"type:varchar(50)"`  // Max 50 characters
	CompanyID      uint                   `gorm:"not null;index"`    // Foreign key to the Category table
	Company        *model_company.Company `gorm:"foreignKey:CompanyID"`
	VerifiedAt     time.Time              `gorm:"type:timestamp"`
	VerifiedMethod string                 `gorm:"type:varchar(50)"`
}

func (UserUpdateModel) TableName() string {
	return "users" // Replace this with your desired table name
}

// Implement the GetCompany method for Product
func (u *User) GetCompany() *model_company.Company {
	return u.Company
}

func FilterByCompanyID(companyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("company_id = ?", companyID)
	}
}
