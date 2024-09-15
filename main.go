// File: main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/utils"
)

func main() {

	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en" // Default to English
	}

	// Initialize the database
	database.InitDB()

	// migrations.MigrateUser()

	user := models.User{Name: "John Doe", Desc: "tree"}

	// Validate and get translated error messages
	errors := utils.ValidateStruct(user)

	// Print the translated error messages
	for _, err := range errors {
		fmt.Printf("%s: %s\n", err.Field, err.Message)
	}

	if err := actions.CreateUser(database.DB, &user); err != nil {
		log.Fatal("Failed to create user:", err)
	}

	// if err != nil {
	// 	fmt.Println("Validation failed:", err)
	// } else {
	// 	fmt.Println("Validation succeeded")
	// }

}
