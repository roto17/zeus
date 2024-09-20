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
	migrations.MigrateLog()

	// Create a new user
	user := models.User{Name: "John Doe yy", Desc: "okokok", Jam: "uuuuu"}
	if err := actions.CreateUser(&user); err != nil {
		log.Fatal("Failed to create user:", err)
	}

}
