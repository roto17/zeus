package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/logs"
)

func MigrateLog() {
	err := database.DB.AutoMigrate(&models.Log{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Log migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "Log migration successful!")
	// log.Println("Log migration successful!")
}
