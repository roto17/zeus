package actions

import (
	"fmt"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	encryptions "github.com/roto17/zeus/lib/encryption"
	"github.com/skip2/go-qrcode"

	model_product_category "github.com/roto17/zeus/lib/models/productcategories"
	model_product "github.com/roto17/zeus/lib/models/products"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

// AddProduct handles the creation of a new product
func AddProduct(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var productEncrypted model_product.ProductEncrypted

	var product model_product.Product

	// Bind the incoming JSON to the product struct
	if err := c.ShouldBindJSON(&productEncrypted); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requestedLanguage)})
		return
	}

	product, ok := encryptions.DecryptObjectID(productEncrypted, &product).(model_product.Product)
	if !ok {
		panic("failed to assert type to Product")
	}

	// Find the user by username
	var searched_category model_product_category.ProductCategory
	if err := db.
		Scopes(model_product_category.
			FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		Where("product_categories.id = ?", product.CategoryID).
		First(&searched_category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("category_not_found", "", requestedLanguage)})
		return
	}

	product.UserID = utils.GetParamIDFromGinClaims(c, "user_id")

	// Validate the incoming product data
	validationErrors := utils.FieldValidationAll(product, requestedLanguage)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the product in the database
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_addition", "", requestedLanguage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("added_successfuly", "", requestedLanguage)})
}

func UpdateProduct(c *gin.Context) {
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var productEncrypted model_product.ProductEncrypted

	var product model_product.Product

	// Bind the incoming JSON to the product struct
	if err := c.ShouldBindJSON(&productEncrypted); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requestedLanguage)})
		return
	}

	product, ok := encryptions.DecryptObjectID(productEncrypted, &product).(model_product.Product)
	if !ok {
		panic("failed to assert type to Product")
	}

	var searched_product model_product.Product
	if err := database.DB.
		Scopes(model_product.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		Where("products.id = ?", product.ID).
		First(&searched_product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("not_found", "", requestedLanguage)})
		return
	}

	// Find the user by username
	var searched_category model_product_category.ProductCategory
	if err := database.DB.
		Scopes(model_product_category.
			FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		Where("product_categories.id = ?", product.CategoryID).
		First(&searched_category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("category_not_found", "", requestedLanguage)})
		return
	}

	product.UserID = utils.GetParamIDFromGinClaims(c, "user_id")

	// Validate the incoming product data
	validationErrors := utils.FieldValidationAll(product, requestedLanguage)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the updated category to the database
	if err := db.
		Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_update", "", requestedLanguage)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("updated_successfully", "", requestedLanguage)})
}

// DeleteProduct deletes a product category by ID
func DeleteProduct(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB

	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	// Check if the category exists
	var product model_product.Product
	if err := db.
		Scopes(model_product.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		First(&product, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Delete the category from the database
	if err := db.Delete(&product).Error; err != nil {

		if strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("fk_issue", "", requested_language)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_deletion", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("delete_successfully", "", requested_language)})
}

// ViewProduct handler
func ViewProduct(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	var product model_product.ProductResponse

	result := database.DB.
		Scopes(model_product.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id"))).
		Preload("Category").
		Preload("User.Company").
		First(&product, idParam)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	encryptedProduct := encryptions.EncryptObjectID(product)

	c.JSON(http.StatusOK, encryptedProduct)
}

// AllProducts handler
func AllProducts(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	// Get pagination parameters from query
	limit, offset := utils.GetPaginationParams(c)

	// Get search query from query parameters
	search := c.DefaultQuery("search", "")

	var products []model_product.ProductResponse
	var totalProducts int64

	// Build base query with search filter
	query := database.DB.Model(&model_product.ProductResponse{})
	query = query.
		Scopes(model_product.FilterByCompanyID(utils.GetParamIDFromGinClaims(c, "company_id")))
	if search != "" {
		query = query.Where("description ILIKE ?", "%"+search+"%") // Case-insensitive search
	}

	// Count total products
	query.Count(&totalProducts)

	// Fetch products with pagination
	result := query.
		Preload("Category").
		Preload("User.Company").
		Limit(limit).
		Offset(offset).
		Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Encrypt product IDs
	encryptedProducts := make([]interface{}, len(products))
	for i, product := range products {
		encryptedProduct := encryptions.EncryptObjectID(product)
		encryptedProducts[i] = encryptedProduct
	}

	// Generate pagination metadata
	pagination := utils.GetPaginationMetadata(c, totalProducts, limit)

	// Return paginated results
	c.JSON(http.StatusOK, gin.H{
		"pagination": pagination,
		"data":       encryptedProducts,
	})
}

func SaveQR(c *gin.Context) {

	content := "wwwwww"

	fmt.Printf("----------------------------\n")
	fmt.Printf("%v", content)
	fmt.Printf("----------------------------\n")

	// if err != nil {
	// 	fmt.Println("Error Encrypting ID Product:", err)
	// }

	// Generate QR code
	qrCode, err := qrcode.New(fmt.Sprintf("%d", content), qrcode.Medium)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	// Set QR code size
	qrCode.DisableBorder = true
	qrCodeImg := qrCode.Image(256) // Generates a 256x256 image

	// Save the QR code image
	outputFile, err := os.Create("lib/images/qrcodes/qrcode.png")
	if err != nil {
		fmt.Println("Error saving QR code:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, qrCodeImg)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}

	fmt.Println("QR code saved to lib/images/qrcodes/qrcode.png")
}
