package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"

	model_product_category "github.com/roto17/zeus/lib/models/productcategories"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

// Add Category
func AddProductCategory(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var category model_product_category.ProductCategory

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	newCategory := model_product_category.ProductCategory{
		Description: category.Description,
		// CreatedAt:   time.Now(),
		// UpdatedAt:   time.Now(),
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(newCategory, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	validatedCategory := model_product_category.ProductCategory{
		Description: newCategory.Description,
		CreatedAt:   newCategory.CreatedAt,
		UpdatedAt:   newCategory.UpdatedAt,
	}

	// Save the user in the database
	if err := db.Create(&validatedCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("registration_failed", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("registration_completed", "", requested_language)})

}
