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
func AddProductToOrder(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var order model_product.Order

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	// Step 2: Get all the ProductIDs from the incoming OrderProducts
	var productIDs []uint
	for _, orderProduct := range order.OrderProducts {
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
	for i := range order.OrderProducts {
		product, exists := productMap[order.OrderProducts[i].ProductID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		order.OrderProducts[i].Product = product
		if product.Weight > 0 {
			order.OrderProducts[i].Price = int64(product.BuyingPrice * product.Weight)
		} else {
			order.OrderProducts[i].Price = int64(product.SellingPrice)
		}

	}

	c.JSON(http.StatusOK, gin.H{"entry": order})

	// fmt.Printf("\n-----%v------------\n", order.OrderProducts)

	// result := db.Create(&order) // `order` should be an Order struct with OrderProducts set

	// if result.Error != nil {
	// 	fmt.Println("Error saving order:", result.Error)
	// 	return
	// }

	result := db.
		Preload("OrderProducts.Product").
		First(&order, order.ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	OrderProducts := order.OrderProducts

	// c.JSON(http.StatusOK, gin.H{"data": order})

	// fmt.Printf("\n++++++__%v___\n", order)

	// OrderProducts := []model_product.OrderProduct{
	// 	{OrderID: order.ID, ProductID: 29, Quantity: 10, Price: 3},
	// 	{OrderID: order.ID, ProductID: 24, Quantity: 2, Price: 2},
	// }

	// for i := range OrderProducts {

	// 	order.OrderProducts[i].Price = int64(order.OrderProducts[i].Product.SellingPrice)

	// 	// order.OrderProducts[i].Price = 100
	// }

	// fmt.Printf("\n----%v---\n", OrderProducts)

	order.OrderProducts = OrderProducts

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

	// // Save order to database
	// if result := db.Create(&order); result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requested_language)})
	// 	return
	// }

	// db.Preload("Product").Find(&order.OrderProducts)

	// db.Preload("Product").Where("order_id = ?", order.ID).Find(&order.OrderProducts)

	c.JSON(http.StatusOK, gin.H{"data": order})

	// c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})
}
