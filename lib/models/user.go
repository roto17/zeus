package models

// User struct with validation
// type User struct {
// 	ID    uint   `gorm:"primaryKey"`                                                 // No validation needed for ID
// 	Name  string `gorm:"type:varchar(255);unique" validate:"required,min=3,max=255"` // Name is required, must be between 3 and 255 characters
// 	Desc  string `gorm:"type:varchar(255);unique" validate:"required"`               // Desc is required
// 	Jam   string `gorm:"type:varchar(255)" validate:"required"`                      // Jam is required
// 	Email string `gorm:"type:varchar(255);unique" validate:"required,email"`         // Jam is required
// }

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"type:varchar(255);unique" validate:"required"`
	Password string `gorm:"type:varchar(50)" validate:"required"`
	Role     string `gorm:"type:varchar(50)" validate:"required"`
}
