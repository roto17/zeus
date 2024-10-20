package main

import (
	"github.com/roto17/zeus/lib/router"

	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
)

func main() {

	// Load Config
	config.LoadConfig()

	// Initialize and load the database
	database.InitDB()

	// Run the server
	router.InitRouter().Run(config.GetEnv("app_running_port")) // Default is :8080
}
