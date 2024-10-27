package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/products"
)

func MigrateProducts() {
	err := database.DB.AutoMigrate(&models.Product{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Product migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "Product migration succeeded!")
	// log.Println("Log migration successful!")
}
