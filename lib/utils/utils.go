package utils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

// Create a validator instance
var validate = validator.New()

// ValidateUser function to validate a User instance
func ValidateStruct(model interface{}) error {
	return validate.Struct(model)
}

// UniqueFieldValidator checks if a specific field is unique across the table.
// It works for any model by accepting the model, field name, and field value.
func UniqueFieldValidator(db *gorm.DB, model interface{}, fieldName string, fieldValue string) error {
	var count int64
	// Dynamically build the query based on the provided field name and value
	query := fmt.Sprintf("%s = ?", fieldName)
	if err := db.Model(model).Where(query, fieldValue).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("%s already exists", fieldName))
	}
	return nil
}
