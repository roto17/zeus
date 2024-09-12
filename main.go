// File: main.go
package main

import (
	"fmt"
	"log"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
)

func main() {

	// Get the GORM DB connection
	db, err := database.DBConnection()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Create a new user
	user := database.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "securepassword",
	}

	// Save the user to the database
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatal("Error creating user:", result.Error)
	} else {
		fmt.Println("User created with ID:", user.ID)
	}

	// Query the user
	var fetchedUser database.User
	if err := db.First(&fetchedUser, "username = ?", "john_doe").Error; err != nil {
		log.Fatal("Error fetching user:", err)
	} else {
		fmt.Printf("Fetched User: %+v\n", fetchedUser)
	}

}
