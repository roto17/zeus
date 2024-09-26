package main

import (
	"log"
	"os"

	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/router"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
)

func main() {
	// Setup PostgreSQL connection
	config.LoadConfig()

	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en" // Default to English
	}

	// Initialize the database
	database.InitDB()

	// Auto-migrate the models
	err := database.DB.AutoMigrate(&models.User{}, &models.Token{})
	if err != nil {
		log.Fatal("failed to migrate the database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Middleware to set the database in the context
	r.Use(func(c *gin.Context) {
		c.Set("db", database.DB)
		c.Next()
	})

	// Routes
	r.POST("/login", router.Login)
	r.POST("/register", router.Register)

	// Run the server
	r.Run() // Default is :8080
}
