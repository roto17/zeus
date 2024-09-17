package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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
func ValidateStruct(model interface{}) ([]ValidationError, error) {
	err := validate.Struct(model)
	var errors []ValidationError

	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			// Get the translated error message

			message := translation.GetTranslation(fieldErr.Tag(), fieldErr.StructField(), "en")

			// Append the custom error message to the errors slice
			errors = append(errors, ValidationError{
				Field:   fieldErr.StructField(),
				Message: message,
			})
		}
	}

	if len(errors) > 0 {
		return errors, nil
	}
	return nil, nil
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

func UniqueFieldValidator_test(db *gorm.DB, model interface{}) ([]ValidationError, error) {
	var uniqueFields []ValidationError
	uniqueFields, err := ValidateStruct(model)

	if err != nil {
		fmt.Printf(err.Error())
	}

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer to get the actual value
	}
	typ := val.Type()

	// Loop through all fields in the struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Check if the field tag contains "unique"
		if strings.Contains(string(field.Tag), "unique") {
			var count int64
			query := fmt.Sprintf("%s = ?", strings.ToLower(fmt.Sprintf("\"%s\"", field.Name)))

			// Check for uniqueness in the database
			if err := db.Model(model).Where(query, value.Interface()).Count(&count).Error; err != nil {
				return nil, err // Return an error if the query fails
			}

			if count > 0 {
				// Field value is not unique
				uniqueFields = append(uniqueFields, ValidationError{
					Field:   field.Name,
					Message: translation.GetTranslation("unique", field.Name, "es"),
				})

			}
		}
	}

	// Return any validation errors found or nil if no errors
	if len(uniqueFields) > 0 {
		return uniqueFields, nil
	}

	return nil, nil
}
