package models

// User struct with validation
type User struct {
	ID    uint   `gorm:"primaryKey"`                                                 // No validation needed for ID
	Name  string `gorm:"type:varchar(255);unique" validate:"required,min=3,max=255"` // Name is required, must be between 3 and 255 characters
	Desc  string `gorm:"type:varchar(255);unique" validate:"required"`               // Desc is required
	Jam   string `gorm:"type:varchar(255)" validate:"required"`                      // Jam is required
	Email string `gorm:"type:varchar(255);unique" validate:"email"`                  // Jam is required
}

// // ValidateUserFields checks both field-level validation and uniqueness
// func ValidateUserFields(db *gorm.DB, user *User) error {
// 	// Field validation using go-playground validator
// 	if err := utils.ValidateStruct(user); err != nil {
// 		return err
// 	}

// 	// Check uniqueness for multiple fields using the reusable unique validator
// 	if err := utils.UniqueFieldValidator(db, &User{}, "name", user.Name); err != nil {
// 		return err
// 	}
// 	if err := utils.UniqueFieldValidator(db, &User{}, "desc", user.Desc); err != nil {
// 		return err
// 	}
// 	if err := utils.UniqueFieldValidator(db, &User{}, "jam", user.Jam); err != nil {
// 		return err
// 	}

// 	return nil
// }
