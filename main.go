// File: main.go
package main

import (
	"fmt"
	"os"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/translation"
)

func main() {

	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en" // Default to English
	}

	// Initialize the database
	database.InitDB()

	// migrations.MigrateUser()

	user := models.User{Name: "John Doe"}

	// Validate and get translated error messages
	errors := translation.ValidateAndTranslate(user, lang)

	// Print the translated error messages
	for _, err := range errors {
		fmt.Printf("%s: %s\n", err.Field, err.Message)
	}

	// if err != nil {
	// 	fmt.Println("Validation failed:", err)
	// } else {
	// 	fmt.Println("Validation succeeded")
	// }

}
