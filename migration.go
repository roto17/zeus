// File: main.go
package main

import (
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/migrations"
)

func main() {

	config.LoadConfig()
	database.InitDB()

	migrations.MigrateNotification()
	migrations.MigrateUser()
	migrations.MigrateToken()

	// Create a new user
	// user := models.User{Name: "John Doe yy", Desc: "okokok", Jam: "uuuuu"}
	// if err := actions.CreateUser(&user); err != nil {
	// 	log.Fatal("Failed to create user:", err)
	// }

}
