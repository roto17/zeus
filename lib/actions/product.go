package actions

import (
	"fmt"
	"image/png"
	"net/http"
	"os"

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

// ViewUser handler
func ViewProduct(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	var product model_product.Product

	result := database.DB.Preload("Category").First(&product, idParam)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	encryptedProduct := encryptions.EncryptObjectID(product)

	// decryptedUser := encryptions.DecryptObjectID(encryptedProduct, &model_product.Product{}).(model_product.Product)

	// Return the encryptedProduct
	c.JSON(http.StatusOK, encryptedProduct)
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
