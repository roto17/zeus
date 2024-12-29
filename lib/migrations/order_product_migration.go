package migrations

import (
	"fmt"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
	"github.com/roto17/zeus/lib/logs"
	models "github.com/roto17/zeus/lib/models/products"
)

func MigrateOrdersProducts() {
	err := database.DB.AutoMigrate(&models.OrderProduct{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("OrderProduct migration failed: %v", err))
		// log.Fatal("Log migration failed:", err)
	}
	logs.AddLog("Info", "roto", "OrderProduct migration succeeded!")
	// log.Println("Log migration successful!")
}
