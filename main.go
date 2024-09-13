// File: main.go
package main

import (
	"log"

	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/migrations"
	"github.com/roto17/zeus/lib/models"
)

func main() {

	// Initialize the database
	database.InitDB()

	migrations.MigrateUser()

	// Create a new user
	user := models.User{Name: "John Doe", Desc: "okok", Jam: "uuu"}
	if err := actions.CreateUser(database.DB, &user); err != nil {
		log.Fatal("Failed to create user:", err)
	}

}
