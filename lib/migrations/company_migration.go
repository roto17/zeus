package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/companies"
)

func MigrateCompany() {
	err := database.DB.AutoMigrate(&models.Company{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Company migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "Company migration succeeded!")
	// log.Println("Log migration successful!")
}
