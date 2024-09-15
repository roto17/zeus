// File: main.go
package main

import (
	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/migrations"
)

func main() {

	// Initialize the database
	database.InitDB()

	migrations.MigrateUser()

	// // Create a new user
	// user := models.User{Name: "John Doe", Desc: "okok", Jam: "uuu"}
	// if err := actions.CreateUser(database.DB, &user); err != nil {
	// 	log.Fatal("Failed to create user:", err)
	// }

}
