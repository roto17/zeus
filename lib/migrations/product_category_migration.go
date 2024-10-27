package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/productcategories"
)

func MigrateProductCategory() {
	err := database.DB.AutoMigrate(&models.ProductCategory{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("ProductCategory migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "ProductCategory migration succeeded!")
	// log.Println("Log migration successful!")
}
