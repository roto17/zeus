package translation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Translation maps for different languages
var translationMap = map[string]map[string]string{
	"en": {
		"required": "{Field} is required",
		"email":    "{Field} must be a valid email",
	},
	"es": {
		"required": "{Field} es obligatorio",
		"email":    "{Field} debe ser un correo electrónico válido",
	},
}

type ValidationError struct {
	Field   string
	Message string
}

// getTranslation fetches the translated message based on the tag and language
func getTranslation(tag string, field string, lang string) string {
	if translationMap[lang] == nil {
		// Fallback to English if the language is not supported
		lang = "en"
	}

	translation := translationMap[lang][tag]
	if translation == "" {
		// Fallback if no specific message for the tag is found
		translation = translationMap["en"][tag]
	}

	// Replace {Field} with the actual field name
	return strings.Replace(translation, "{Field}", field, 1)
}

// validateAndTranslate runs validation and translates error messages
func ValidateAndTranslate(model interface{}, lang string) []ValidationError {
	validate := validator.New()
	err := validate.Struct(model)
	var errors []ValidationError

	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			// Get the translated error message
			message := getTranslation(fieldErr.Tag(), fieldErr.StructField(), lang)

			// Append the custom error message to the errors slice
			errors = append(errors, ValidationError{
				Field:   fieldErr.StructField(),
				Message: message,
			})
		}
	}
	return errors
}
