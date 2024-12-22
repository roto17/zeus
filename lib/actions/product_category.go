package actions

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	encryptions "github.com/roto17/zeus/lib/encryption"
	model_product_category "github.com/roto17/zeus/lib/models/productcategories"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
	"gorm.io/gorm"
)

// Add Category
func AddProductCategory(c *gin.Context) {

	requestedLanguage := utils.GetHeaderVarToString(c.Get("requestedLanguage"))
	db := database.DB
	var category model_product_category.ProductCategory

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requestedLanguage)})
		return
	}

	category.UserID = utils.GetParamIDFromGinClaims(c, "user_id")

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(category, requestedLanguage)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the user in the database
	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requestedLanguage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requestedLanguage)})
}

// Update Category
func UpdateProductCategory(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requestedLanguage"))
	db := database.DB
	var encryptedCategory model_product_category.ProductCategoryEncrypted
	var category model_product_category.ProductCategory

	// Bind the incoming JSON to the category struct
	if err := c.ShouldBindJSON(&encryptedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requestedLanguage)})
		return
	}

	category, ok := encryptions.DecryptObjectID(encryptedCategory, &category).(model_product_category.ProductCategory)
	if !ok {
		panic("failed to assert type to ProductCategory")
	}

	// Check if the category exists by ID
	if err := db.
		Scopes(model_product_category.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		Where("product_categories.id = ?", category.ID).
		First(&model_product_category.ProductCategory{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requestedLanguage)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("internal_error", "", requestedLanguage)})
		}
		return
	}

	// Update fields from the input
	category.UserID = utils.GetParamIDFromGinClaims(c, "user_id")

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(category, requestedLanguage)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the updated category to the database
	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_update", "", requestedLanguage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("updated_successfully", "", requestedLanguage)})
}

// DeleteProductCategory deletes a product category by ID
func DeleteProductCategory(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requestedLanguage"))
	db := database.DB

	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	// Check if the category exists
	var category model_product_category.ProductCategory
	if err := db.
		Scopes(model_product_category.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		First(&category, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requestedLanguage)})
		return
	}

	// Delete the category from the database
	if err := db.Delete(&category).Error; err != nil {

		if strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("fk_issue", "", requestedLanguage)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_deletion", "", requestedLanguage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("delete_successfully", "", requestedLanguage)})
}

// ViewUser handler
func ViewProductCategory(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requestedLanguage"))
	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	var category model_product_category.ProductCategory

	result := database.DB.
		Scopes(model_product_category.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		Preload("User.Company").
		First(&category, idParam)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requestedLanguage)})
		return
	}

	encryptedCategory := encryptions.EncryptObjectID(category)

	// // Return the encryptedUser
	c.JSON(http.StatusOK, encryptedCategory)
}

// AllProductCategories handler
func AllProductCategories(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requestedLanguage"))

	// Get pagination parameters from query
	limit, offset := utils.GetPaginationParams(c)

	// Get search query from query parameters
	search := c.DefaultQuery("search", "")

	var categories []model_product_category.ProductCategory
	var totalCategories int64

	// Build base query with search filter
	query := database.DB.Model(&model_product_category.ProductCategory{})
	query = query.Preload("User.Company").
		Scopes(model_product_category.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id")))

	if search != "" {
		query = query.
			// Preload("User").
			Where("description ILIKE ?", "%"+search+"%") // Case-insensitive search for category names
	}

	// Count total categories
	query.Count(&totalCategories)

	// Fetch categories with pagination
	result := query.Limit(limit).Offset(offset).Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requestedLanguage)})
		return
	}

	// Encrypt category IDs
	encryptedCategories := make([]interface{}, len(categories))
	for i, category := range categories {
		encryptedCategory := encryptions.EncryptObjectID(category)
		encryptedCategories[i] = encryptedCategory
	}

	// Generate pagination metadata
	pagination := utils.GetPaginationMetadata(c, totalCategories, limit)

	// Return paginated results
	c.JSON(http.StatusOK, gin.H{
		"pagination": pagination,
		"data":       encryptedCategories,
	})
}
