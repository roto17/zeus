package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/models"
)

func MigrateLog() {
	err := database.DB.AutoMigrate(&models.Log{})
	if err != nil {
		actions.Fatal("roto", fmt.Sprintf("Log migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	actions.Info("roto", "Log migration successful!")
	// log.Println("Log migration successful!")
}
