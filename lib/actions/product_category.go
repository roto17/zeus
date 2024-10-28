package actions

import (
	"net/http"
	"strconv"

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
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(newCategory, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	validatedCategory := model_product_category.ProductCategory{
		Description: newCategory.Description,
	}

	// Save the user in the database
	if err := db.Create(&validatedCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})
}

// Update Category
func UpdateProductCategory(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var category model_product_category.ProductCategory

	// Bind the incoming JSON to the category struct
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	// Fetch the existing category by ID
	var existingCategory model_product_category.ProductCategory
	if err := db.First(&existingCategory, category.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("category_not_found", "", requested_language)})
		return
	}

	// Update fields from the input
	existingCategory.Description = category.Description

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(existingCategory, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the updated category to the database
	if err := db.Save(&existingCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_update", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("updated_successfully", "", requested_language)})
}

// DeleteProductCategory deletes a product category by ID
func DeleteProductCategory(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB

	// Get the category ID from the URL parameter
	categoryID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_id", "", requested_language)})
		return
	}

	// Check if the category exists
	var category model_product_category.ProductCategory
	if err := db.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("category_not_found", "", requested_language)})
		return
	}

	// Delete the category from the database
	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_delete", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("delete_successfully", "", requested_language)})
}
