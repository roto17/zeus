package migrations

import (
	"log"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/models"
)

func MigrateUser() {
	err := database.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("User migration failed:", err)
	}
	log.Println("User migration successful!")
}
