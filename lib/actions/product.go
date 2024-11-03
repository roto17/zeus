package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"

	model_product_category "github.com/roto17/zeus/lib/models/productcategories"
	model_product "github.com/roto17/zeus/lib/models/products"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

// AddProduct handles the creation of a new product
func AddProduct(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var productInput model_product.ProductInput

	// Bind the incoming JSON to the product struct
	if err := c.ShouldBindJSON(&productInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requestedLanguage)})
		return
	}

	// Find the user by username
	var searched_category model_product_category.ProductCategory
	if err := database.DB.Where("id = ?", productInput.CategoryID).First(&searched_category).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No mathcing Catgeory"})
		return
	}

	productValidation := model_product.Product{
		Description: productInput.Description,
		QRCode:      productInput.QRCode,
		CategoryID:  productInput.CategoryID,
		Category:    searched_category,
	}

	// Validate the incoming product data
	validationErrors := utils.FieldValidationAll(productValidation, requestedLanguage)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Ensure CategoryID is set
	if productInput.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("category_not_found", "", requestedLanguage)})
		return
	}

	// Save the product in the database
	if err := db.Create(&productValidation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requestedLanguage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requestedLanguage)})
}
