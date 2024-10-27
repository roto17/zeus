package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/users"
)

func MigrateUser() {
	err := database.DB.AutoMigrate(&models.User{})
	if err != nil {
		// log.Fatal("User migration failed:", err)
		logs.AddLog("Fatal", "roto", fmt.Sprintf("User migration failed: %v", err))
	}
	// log.Println("User migration successful!")
	logs.AddLog("Info", "roto", "User migration succeeded!")
}
