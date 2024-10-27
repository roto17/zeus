package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/notifications"
)

func MigrateNotification() {
	err := database.DB.AutoMigrate(&models.Notification{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Notification migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "Notification migration succeeded!")
	// log.Println("Log migration successful!")
}
