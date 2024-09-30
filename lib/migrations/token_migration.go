package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/tokens"
)

func MigrateToken() {
	err := database.DB.AutoMigrate(&models.Token{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Token migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "Token migration successful!")
	// log.Println("Log migration successful!")
}
