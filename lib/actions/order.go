package actions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	"gorm.io/gorm/clause"

	// encryptions "github.com/roto17/zeus/lib/encryption"
	// "gorm.io/gorm"

	model_product "github.com/roto17/zeus/lib/models/products"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

// Add Category
func AddOrder(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var order model_product.Order

	// // Bind the incoming JSON to the user struct
	// if err := c.ShouldBindJSON(&order); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
	// 	return
	// }

	order = model_product.Order{
		UserID:    utils.GetParamIDFromGinClaims(c, "user_id"),
		Direction: "+",
		Status:    "pending",
	}

	result := db.Create(&order) // `order` should be an Order struct with OrderProducts set

	if result.Error != nil {
		fmt.Println("Error saving order:", result.Error)
		return
	}

	fmt.Printf("\n++++++__%v___\n", order)

	OrderProducts := []model_product.OrderProduct{
		{OrderID: order.ID, ProductID: 29, Quantity: 5},
		{OrderID: order.ID, ProductID: 24, Quantity: 3},
	}

	order.OrderProducts = OrderProducts

	result3 := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "order_id"},   // Specify the column `order_id` in the conflict clause (from `order_products` table).
			{Name: "product_id"}, // Specify the column `product_id` in the conflict clause (from `order_products` table).
		},
		DoUpdates: clause.AssignmentColumns([]string{"quantity"}), // Update the `quantity` field if there is a conflict.
	}).Create(&order.OrderProducts) // Insert the order products and handle conflicts.

	if result3.Error != nil {
		fmt.Println("Error:", result3.Error)
	} else {
		fmt.Println("Product added to order")
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})
}
