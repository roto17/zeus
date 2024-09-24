package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/logs"
	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/translation"
)

// Create a validator instance
var validate = validator.New()

// ValidateUser function to validate a User instance
func ValidateStruct(model interface{}, language string) []models.ValidationError {
	err := validate.Struct(model)
	var errors []models.ValidationError

	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			// Get the translated error message

			message := translation.GetTranslation(fieldErr.Tag(), fieldErr.StructField(), language)

			// Append the custom error message to the errors slice
			errors = append(errors, models.ValidationError{
				Field:   fieldErr.StructField(),
				Message: message,
			})
		}
	}

	return errors
}

func FieldValidationAll(model interface{}, language string) []models.ValidationError {

	listOfErrors := ValidateStruct(model, language)

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
			if err := database.DB.Model(model).Where(query, value.Interface()).Count(&count).Error; err != nil {
				// return nil, err // Return an error if the query fails
				logs.AddLog("Fatal", "roto", fmt.Sprintf("Failed to connect to the database:%s", err))
			}

			if count > 0 {
				// Field value is not unique
				listOfErrors = append(listOfErrors, models.ValidationError{
					Field:   field.Name,
					Message: translation.GetTranslation("unique", field.Name, language),
				})

			}
		}
	}

	// Return any validation errors found or nil if no errors
	// if len(listOfErrors) > 0 {
	return listOfErrors
	// }

	// return nil
}
