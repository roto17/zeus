package main

import (
	"github.com/roto17/zeus/lib/router"

	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
)

func main() {
	// Setup PostgreSQL connection
	config.LoadConfig()

	// lang := os.Getenv("LANG")
	// if lang == "" {
	// 	lang = "en" // Default to English
	// }

	// Initialize the database
	database.InitDB()

	// Auto-migrate the models
	// err := database.DB.AutoMigrate(&models.User{}, &models.Token{})
	// if err != nil {
	// 	log.Fatal("failed to migrate the database:", err)
	// }

	// Initialize Gin router

	// Run the server
	router.InitRouter().Run(config.GetEnv("app_running_port")) // Default is :8080
}
