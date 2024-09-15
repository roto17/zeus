package utils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"

	"github.com/roto17/zeus/lib/translation"
)

type ValidationError struct {
	Field   string
	Message string
}

// Create a validator instance
var validate = validator.New()

// ValidateUser function to validate a User instance
func ValidateStruct(model interface{}) []ValidationError {
	err := validate.Struct(model)
	var errors []ValidationError

	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			// Get the translated error message

			message := translation.GetTranslation(fieldErr.Tag(), fieldErr.StructField(), "es")

			// Append the custom error message to the errors slice
			errors = append(errors, ValidationError{
				Field:   fieldErr.StructField(),
				Message: message,
			})
		}
	}
	return errors
	// return validate.Struct(model)
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
