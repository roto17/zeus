package actions

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/roto17/zeus/lib/database"
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

func SaveQR(c *gin.Context) {

	content := "https://example.com"

	// Generate QR code
	qrCode, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	// Set QR code size
	qrCode.DisableBorder = true
	qrCodeImg := qrCode.Image(256) // Generates a 256x256 image

	// Open the logo image
	logoFile, err := os.Open("lib/images/logo/logo.png")
	if err != nil {
		fmt.Println("Error opening logo:", err)
		return
	}
	defer logoFile.Close()

	// Decode logo image
	logoImg, _, err := image.Decode(logoFile)
	if err != nil {
		fmt.Println("Error decoding logo:", err)
		return
	}

	// Resize the logo to be smaller
	logoSize := uint(50) // Set desired logo size
	resizedLogo := resize.Resize(logoSize, logoSize, logoImg, resize.Lanczos3)

	// Position the logo in the bottom-right corner
	offset := image.Pt(qrCodeImg.Bounds().Dx()-resizedLogo.Bounds().Dx()-10, qrCodeImg.Bounds().Dy()-resizedLogo.Bounds().Dy()-10)

	// Create a new RGBA image with the same size as the QR code
	finalImg := image.NewRGBA(qrCodeImg.Bounds())

	// Draw the QR code on the final image
	draw.Draw(finalImg, qrCodeImg.Bounds(), qrCodeImg, image.Point{}, draw.Src)

	// Overlay the resized logo in the bottom-right corner of the QR code
	draw.Draw(finalImg, resizedLogo.Bounds().Add(offset), resizedLogo, image.Point{}, draw.Over)

	// Save the final image
	outputFile, err := os.Create("lib/images/qrcodes/qrcode_with_bottom_corner_logo.png")
	if err != nil {
		fmt.Println("Error saving QR code:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, finalImg)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}

	fmt.Println("QR code with logo saved to lib/images/qrcodes/qrcode_with_bottom_corner_logo.png")
}
