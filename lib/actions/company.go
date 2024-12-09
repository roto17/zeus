package actions

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	encryptions "github.com/roto17/zeus/lib/encryption"

	model_company "github.com/roto17/zeus/lib/models/companies"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

// Add Category
func AddCompany(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var company model_company.Company

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	newCompany := model_company.Company{
		Description: company.Description,
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(newCompany, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the user in the database
	if err := db.Create(&newCompany).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requested_language)})
}

// ViewUser handler
func ViewCompany(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	var company model_company.Company

	result := database.DB.First(&company, idParam)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	encryptedCompany := encryptions.EncryptObjectID(company)

	// // Return the encryptedUser
	c.JSON(http.StatusOK, encryptedCompany)
}

// AllProductCategories handler
func AllCompanies(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	// Get pagination parameters from query
	limit, offset := utils.GetPaginationParams(c)

	// Get search query from query parameters
	search := c.DefaultQuery("search", "")

	var companies []model_company.Company
	var totalCompanies int64

	// Build base query with search filter
	query := database.DB.Model(&model_company.Company{})
	if search != "" {
		query = query.Where("description ILIKE ?", "%"+search+"%") // Case-insensitive search for company names
	}

	// Count total categories
	query.Count(&totalCompanies)

	// Fetch categories with pagination
	result := query.Limit(limit).Offset(offset).Find(&companies)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Encrypt category IDs
	encryptedCompanies := make([]interface{}, len(companies))
	for i, company := range companies {
		encryptedCompany := encryptions.EncryptObjectID(company)
		encryptedCompanies[i] = encryptedCompany
	}

	pagination := utils.GetPaginationMetadata(c, totalCompanies, limit)

	// Return paginated results
	c.JSON(http.StatusOK, gin.H{
		"pagination": pagination,
		"data":       encryptedCompanies,
	})
}

// Update Category
func UpdateCompany(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var encryptedCompany model_company.CompanyEncrypted
	var company model_company.Company

	// Bind the incoming JSON to the company struct
	if err := c.ShouldBindJSON(&encryptedCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	company, ok := encryptions.DecryptObjectID(encryptedCompany, &company).(model_company.Company)
	if !ok {
		panic("failed to assert type to Company")
	}

	// Fetch the existing company by ID
	var existingCompany model_company.Company
	if err := db.First(&existingCompany, company.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Update fields from the input
	existingCompany.Description = company.Description

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(existingCompany, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the updated category to the database
	if err := db.Save(&existingCompany).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_update", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("updated_successfully", "", requested_language)})
}

// DeleteCompany deletes a company by ID
func DeleteCompany(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB

	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	// Check if the category exists
	var company model_company.Company
	if err := db.First(&company, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Delete the category from the database
	if err := db.Delete(&company).Error; err != nil {

		if strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("fk_issue", "", requested_language)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_deletion", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("delete_successfully", "", requested_language)})
}
