package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"

	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/logs"
	model_validation "github.com/roto17/zeus/lib/models/validations"
	"github.com/roto17/zeus/lib/translation"
)

// Create a validator instance
var validate = validator.New()

// ValidateUser function to validate a User instance
func ValidateStruct(model interface{}, language string) []model_validation.ValidationError {
	err := validate.Struct(model)
	var errors []model_validation.ValidationError

	// fmt.Printf("**************\n")
	// fmt.Printf("%s", model.(model_user.User).Password)
	// fmt.Printf("**************\n")

	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			// Get the translated error message

			message := strings.Replace(translation.GetTranslation(fieldErr.Tag(), fieldErr.StructField(), language), "{Field}", fieldErr.StructField(), 1) // Replace {Field} with the actual field name

			// message = strings.Replace(message, "{Field}", fieldErr.StructField(), 1)

			// Check if the error tag is "oneof"
			if fieldErr.Tag() == "oneof" {
				// Dynamically retrieve the allowed values from the struct tag
				allowedValues := getOneOfTagValue(model, fieldErr.StructField())

				// Replace the {Values} placeholder in the error message
				message = strings.Replace(message, "{values}", allowedValues, 1)
			}
			fmt.Printf("%s", fieldErr.StructField())

			// Append the custom error message to the errors slice
			errors = append(errors, model_validation.ValidationError{
				Field:   fieldErr.StructField(),
				Message: message,
			})
		}
	}

	return errors
}

func FieldValidationAll(model interface{}, language string) []model_validation.ValidationError {

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
				listOfErrors = append(listOfErrors, model_validation.ValidationError{
					Field:   field.Name,
					Message: strings.Replace(translation.GetTranslation("unique", field.Name, language), "{Field}", field.Name, 1), // Replace {Field} with the actual field name
				})

			}
		}
	}

	return listOfErrors
}

// Helper function to extract the "oneof" tag value from the struct
func getOneOfTagValue(model interface{}, fieldName string) string {
	// Get the type of the model
	val := reflect.ValueOf(model)

	// Ensure that we're dealing with a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference the pointer to get to the struct
	}

	// Get the field by name and check if it exists
	field, ok := val.Type().FieldByName(fieldName)
	if !ok {
		return "" // Return an empty string if field not found
	}

	// Extract the `validate` tag and split to find the "oneof" values
	tag := field.Tag.Get("validate")
	if strings.Contains(tag, "oneof=") {
		// Extract the values inside the "oneof" tag
		oneOfValues := strings.Split(tag, "oneof=")[1]
		return oneOfValues
	}

	return ""
}

func SendVerificationEmail(userEmail, token, appBaseURL, smtpUser, smtpPass, smtpHost string, smtpPort int) error {
	// Create the verification URL
	verificationURL := fmt.Sprintf("%s/verify-email?signature=%s", appBaseURL, token)

	// Email content
	subject := "Email Verification"
	// body := fmt.Sprintf("Please click the following link to verify your email: %s", verificationURL)

	// HTML version of the email body
	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<p>Please click the following link to verify your email:</p>
				<a href="%s">Verify Email</a>
			</body>
		</html>`, verificationURL)

	// Set up the email message
	message := gomail.NewMessage()
	message.SetHeader("From", smtpUser)
	message.SetHeader("To", userEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", htmlBody) // plain text
	// message.AddAlternative("text/html", htmlBody) // HTML version

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
