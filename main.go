// File: main.go
package main

import (
	"fmt"
	"os"

	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/utils"
)

func main() {

	config.LoadConfig()

	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en" // Default to English
	}

	// Initialize the database
	database.InitDB()

	// migrations.MigrateUser()

	user := models.User{Name: "John Doe2", Desc: "Desc", Jam: "OKOK", Email: "test@test.com"}

	// Validate and get translated error messages
	validationerrors, error := utils.FieldValidationAll(user, "en")

	if validationerrors != nil {
		for _, err := range validationerrors {
			fmt.Printf("%s: %s\n", err.Field, err.Message)
		}
	}

	if error != nil {
		fmt.Println("Validation failed:", error)
	}

}
