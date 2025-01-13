package actions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/translation"
	"github.com/roto17/zeus/lib/utils"
	"gorm.io/gorm/clause"

	// encryptions "github.com/roto17/zeus/lib/encryption"
	// "gorm.io/gorm"

	model_product "github.com/roto17/zeus/lib/models/products" // Assuming translation package handles translations
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

	// result := db.Create(&order) // `order` should be an Order struct with OrderProducts set

	// if result.Error != nil {
	// 	fmt.Println("Error saving order:", result.Error)
	// 	return
	// }

	// result := db.
	// 	Preload("OrderProducts.Product").
	// 	First(&order, 7)

	// if result.Error != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
	// 	return
	// }

	// fmt.Printf("\n++++++__%v___\n", order)

	// OrderProducts := []model_product.OrderProduct{
	// 	{OrderID: order.ID, ProductID: 29, Quantity: 2, Price: 111},
	// 	{OrderID: order.ID, ProductID: 24, Quantity: 2, Price: 8},
	// }

	// order.OrderProducts = OrderProducts

	// result = db.Clauses(clause.OnConflict{
	// 	Columns: []clause.Column{
	// 		{Name: "order_id"},   // Specify the column `order_id` in the conflict clause (from `order_products` table).
	// 		{Name: "product_id"}, // Specify the column `product_id` in the conflict clause (from `order_products` table).
	// 	},
	// 	DoUpdates: clause.AssignmentColumns([]string{"quantity", "price"}), // Update the `quantity` field if there is a conflict.
	// }).Create(&order.OrderProducts) // Insert the order products and handle conflicts.

	// if result.Error != nil {
	// 	fmt.Println("Error:", result.Error)
	// } else {
	// 	fmt.Println("Product added to order")
	// }

	// encryptedOrder := encryptions.EncryptObjectID(order)

	// Save order to database
	if result := db.Create(&order); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": order.ID})

	// c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})
}

// Add Category
func AddProductToStock(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var orderInput model_product.Order

	var order model_product.Order

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&orderInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	result := db.
		Preload("OrderProducts.Product").
		First(&order, orderInput.ID)

	// result2 := db.Where("order_id = ?", orderInput.ID).Delete(&model_product.OrderProduct{})

	// if result2.Error != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Can't delete order products"})
	// 	return
	// }

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	if order.Status == "pending" { // Step 2: Get all the ProductIDs from the incoming OrderProducts

		if order.Direction == "+" {

			var productIDs []uint
			for _, orderProduct := range orderInput.OrderProducts {
				productIDs = append(productIDs, orderProduct.ProductID)
			}

			// Step 3: Retrieve all associated Products in a single query using `IN`
			var products []model_product.Product
			if err := db.Where("id IN (?)", productIDs).Find(&products).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Products not found"})
				return
			}

			// Step 4: Create a map of ProductID -> Product for fast lookup
			productMap := make(map[uint]model_product.Product)
			for _, product := range products {
				productMap[product.ID] = product
			}

			// Step 5: Associate each OrderProduct with its corresponding Product
			for i := range orderInput.OrderProducts {
				product, exists := productMap[orderInput.OrderProducts[i].ProductID]
				if !exists {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
					return
				}
				orderInput.OrderProducts[i].Product = product
				// if product.Weight > 0 {
				// 	orderInput.OrderProducts[i].Price = int64(product.BuyingPrice * product.Weight)
				// } else {
				// 	orderInput.OrderProducts[i].Price = int64(product.SellingPrice)
				// }

			}

			order.OrderProducts = orderInput.OrderProducts

			result = db.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "order_id"},   // Specify the column `order_id` in the conflict clause (from `order_products` table).
					{Name: "product_id"}, // Specify the column `product_id` in the conflict clause (from `order_products` table).
				},
				DoUpdates: clause.AssignmentColumns([]string{"quantity", "price"}), // Update the `quantity` field if there is a conflict.
			}).Create(&order.OrderProducts) // Insert the order products and handle conflicts.

			if result.Error != nil {
				fmt.Println("Error:", result.Error)
			} else {
				fmt.Println("Product added to order")
			}

			// encryptedOrder := encryptions.EncryptObjectID(order)

			c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})

			// c.JSON(http.StatusOK, gin.H{"data": order})
			return

		} else {

			var productIDs []uint
			for _, orderProduct := range orderInput.OrderProducts {
				productIDs = append(productIDs, orderProduct.ProductID)
			}

			// Step 3: Retrieve all associated Products in a single query using `IN`
			var products []model_product.Product
			if err := db.Where("id IN (?)", productIDs).Find(&products).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Products not found"})
				return
			}

			// Step 4: Create a map of ProductID -> Product for fast lookup
			productMap := make(map[uint]model_product.Product)
			for _, product := range products {
				productMap[product.ID] = product
			}

			var errorProducts []ErrorProduct
			// Step 5: Associate each OrderProduct with its corresponding Product
			for i := range orderInput.OrderProducts {
				product, exists := productMap[orderInput.OrderProducts[i].ProductID]
				if !exists {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
					return
				}
				orderInput.OrderProducts[i].Product = product

				if product.Weight > 0 {
					orderInput.OrderProducts[i].Price = int64(product.BuyingPrice * product.Weight)
				} else {
					orderInput.OrderProducts[i].Price = int64(product.SellingPrice)
				}

				if orderInput.OrderProducts[i].Quantity > product.Quantity {
					errorProducts = append(errorProducts,
						ErrorProduct{
							ProductID: int(product.ID),
							Msg:       fmt.Sprintf("Requested %d of this product : %s ,but we only have %d", orderInput.OrderProducts[i].Quantity, product.Description, product.Quantity),
						})

					// fmt.Println("Requested QTT is higher than stock")
				}
				//  else {
				// 	fmt.Println("Requested QTT exists in stock")
				// }

			}

			if len(errorProducts) > 0 {
				c.JSON(http.StatusOK, gin.H{"error": errorProducts})
				return
			}

			order.OrderProducts = orderInput.OrderProducts

			result = db.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "order_id"},   // Specify the column `order_id` in the conflict clause (from `order_products` table).
					{Name: "product_id"}, // Specify the column `product_id` in the conflict clause (from `order_products` table).
				},
				DoUpdates: clause.AssignmentColumns([]string{"quantity", "price"}), // Update the `quantity` field if there is a conflict.
			}).Create(&order.OrderProducts) // Insert the order products and handle conflicts.

			if result.Error != nil {
				fmt.Println("Error:", result.Error)
			} else {
				fmt.Println("Product added to order")
			}

			c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})
			return

		}

	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Order is closed"})
		return
	}

}

type ErrorProduct struct {
	ProductID int    `json:"product_id"`
	Msg       string `json:"msg"`
}
