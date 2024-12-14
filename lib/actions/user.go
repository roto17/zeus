package actions

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	encryptions "github.com/roto17/zeus/lib/encryption"
	model_company "github.com/roto17/zeus/lib/models/companies"
	model_token "github.com/roto17/zeus/lib/models/tokens"
	model_user "github.com/roto17/zeus/lib/models/users"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

// // Register handles registration
// func VerifyBySMS(c *gin.Context) {

// 	utils.SendSMS()

// 	c.JSON(http.StatusOK, "MSG SENT")

// }

// Register handles registration
func Register(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var encryptedUser model_user.EncryptedUser

	var new_user model_user.User

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&encryptedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	new_user, ok := encryptions.DecryptObjectID(encryptedUser, &new_user).(model_user.User)
	if !ok {
		panic("failed to assert type to User")
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(new_user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("password_hashing_failed", "", requested_language)})
		return
	}

	var searched_company model_company.Company
	if err := database.DB.Where("id = ?", new_user.CompanyID).First(&searched_company).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("company_not_found", "", requested_language)})
		return
	}

	newUser := model_user.User{
		Email:      new_user.Email,
		FirstName:  new_user.FirstName,
		LastName:   new_user.LastName,
		Username:   new_user.Username,
		Password:   hashedPassword,
		Role:       new_user.Role,
		MiddleName: new_user.MiddleName,
		CompanyID:  new_user.CompanyID,
		Company:    &searched_company,
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(newUser, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the user in the database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("registration_failed", "", requested_language)})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("token_generation_failed", "", requested_language)})
		return
	}

	if err := utils.SendVerificationEmail(newUser.Email, token, config.GetEnv("appBaseURL"), config.GetEnv("smtpUser"), config.GetEnv("smtpPass"), config.GetEnv("smtpHost"), utils.StringToInt((config.GetEnv("smtpPort")))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("failed_to_send_verification_email", "", requested_language)})
		return
	}

	// toRoles := []string{"admin", "manager", "user"}
	// notifications.Notify(newUser.Username, newUser.Role, strings.Join(toRoles, ","), "added success!") // Call the RegisterUser function to send a notification

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("registration_completed", "", requested_language)})

}

// Login handles login and JWT generation
func Login(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	var loginData model_user.LoginUserInput

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(loginData, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Find the user by username
	var user model_user.User
	if err := database.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil || !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("password_username_wrong", "", requested_language)})
		return
	}

	// Check if the password matches
	if user.VerifiedAt.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("verify_account", "", requested_language)})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("token_generation_failed", "", requested_language)})
		return
	}

	ip_addr := utils.Coalesce(c.GetHeader("X-Forwarded-For"), c.ClientIP())

	// Save the token and expiration in the database
	newToken := model_token.Token{
		Token:      token,
		UserID:     user.ID,
		User:       user,
		IPAddress:  ip_addr,
		DeviceName: utils.DeviceNameString(c.GetHeader("User-Agent")),
	}

	if err := database.DB.Create(&newToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("token_save_failed", "", requested_language)})
		return
	}

	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Logout handles logout by invalidating the JWT token
func Logout(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("required_header", "", requested_language)})
		c.Abort()
		return
	}

	// Strip "Bearer " if it's included in the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	result := database.DB.Where("token = ?", tokenString).Delete(&model_token.Token{})

	if result.RowsAffected == 0 || result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("cannot_logout", "", requested_language)})
		c.Abort()
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("logout_successful", "", requested_language)})
}

// Logout handles logout by invalidating the JWT token
func LogoutAll(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("required_header", "", requested_language)})
		c.Abort()
		return
	}

	// Strip "Bearer " if it's included in the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.GetEnv("secretkey")), nil
	})

	if err != nil || !token.Valid {

		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_or_expired_token", "", requested_language)})
		c.Abort()
		return
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token_claims", "", requested_language)})
		c.Abort()
		return
	}

	userId := claims["user_id"].(string)
	// userRole := claims["role"].(string)

	result := database.DB.Where("user_id = ?", userId).Delete(&model_token.Token{})

	if result.RowsAffected == 0 || result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("cannot_logout", "", requested_language)})
		c.Abort()
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("logout_successful", "", requested_language)})
}

// ViewProduct handler
func ViewUser(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	var user model_user.User

	result := database.DB.
		Scopes(model_user.FilterByCompanyID(utils.GetCompanyIDFromGinClaims(c))).
		Preload("Company").First(&user, idParam)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// isMatching := utils.IsCompanyIDMatching(&user, utils.GetCompanyIDFromGinClaims(c))

	// if !isMatching {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("Insufficient_permissions", "", requested_language)})
	// 	return
	// }

	encryptedUser := encryptions.EncryptObjectID(user)

	// decryptedUser := encryptions.DecryptObjectID(encryptedProduct, &model_product.Product{}).(model_product.Product)

	// Return the encryptedProduct
	c.JSON(http.StatusOK, encryptedUser)
}

// Update User
func UpdateUser(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var encryptedUser model_user.EncryptedUser
	var user model_user.User

	// Bind the incoming JSON to the category struct
	if err := c.ShouldBindJSON(&encryptedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	user, ok := encryptions.DecryptObjectID(encryptedUser, &user).(model_user.User)
	if !ok {
		panic("failed to assert type to User")
	}

	// Fetch the existing category by ID
	var existingUser model_user.UserUpdateModel
	if err := db.First(&existingUser, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	var searched_company model_company.Company
	if err := database.DB.Where("id = ?", user.CompanyID).First(&searched_company).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("company_not_found", "", requested_language)})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("password_hashing_failed", "", requested_language)})
		return
	}

	// Update fields from the input
	existingUser.FirstName = user.FirstName
	existingUser.MiddleName = user.MiddleName
	existingUser.LastName = user.LastName
	existingUser.Password = hashedPassword
	existingUser.CompanyID = user.CompanyID
	existingUser.Company = &searched_company

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(existingUser, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Save the updated category to the database
	if err := db.Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_update", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("updated_successfully", "", requested_language)})
}

// DeleteUser deletes a product category by ID
func DeleteUser(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB

	escapedID := utils.GetHeaderVarToString(c.Get("escapedID"))

	idParam := utils.DecryptID(escapedID)

	// Check if the user exists
	var user model_user.User
	if err := db.First(&user, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Delete the user from the database
	if err := db.Delete(&user).Error; err != nil {

		if strings.Contains(err.Error(), "23503") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("fk_issue", "", requested_language)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("faild_deletion", "", requested_language)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("delete_successfully", "", requested_language)})
}

// AllUsers handler
func AllUsers(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	// Get pagination parameters from query
	limit, offset := utils.GetPaginationParams(c)

	// Get search query from query parameters
	search := c.DefaultQuery("search", "")

	var users []model_user.User
	var totalUsers int64

	// Build base query with search filter
	query := database.DB.Model(&model_user.User{})
	query = query.Scopes(model_user.FilterByCompanyID(utils.GetCompanyIDFromGinClaims(c)))
	if search != "" {
		query = query.Where("first_name ILIKE ?", "%"+search+"%") // Case-insensitive search for category names
	}

	// Count total categories
	query.Count(&totalUsers)

	// Fetch categories with pagination
	result := query.Limit(limit).Offset(offset).Preload("Company").Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Encrypt category IDs
	encryptedUsers := make([]interface{}, len(users))
	for i, user := range users {
		encryptedUser := encryptions.EncryptObjectID(user)
		encryptedUsers[i] = encryptedUser
	}

	// Generate pagination metadata
	pagination := utils.GetPaginationMetadata(c, totalUsers, limit)

	// Return paginated results
	c.JSON(http.StatusOK, gin.H{
		"pagination": pagination,
		"data":       encryptedUsers,
	})
}

// VerifyByMail handles email verification
func VerifyByMail(c *gin.Context) {
	tokenString := c.Query("signature")
	requestedLanguage := utils.GetHeaderVarToString(c.Get("requested_language"))

	if tokenString == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": translation.GetTranslation("signature_required", "", requestedLanguage)})
		return
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.GetEnv("secretkey")), nil
	})

	if err != nil || !token.Valid {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": translation.GetTranslation("invalid_or_expired_token", "", requestedLanguage)})
		c.Abort()
		return
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": translation.GetTranslation("invalid_token_claims", "", requestedLanguage)})
		c.Abort()
		return
	}

	userID := utils.StringToInt(claims["user_id"].(string))

	var user model_user.User
	result := database.DB.First(&user, userID)

	if result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": translation.GetTranslation("not_found", "", requestedLanguage)})
		c.Abort()
		return
	}

	if !user.VerifiedAt.IsZero() {
		c.HTML(http.StatusOK, "success.html", gin.H{"message": translation.GetTranslation("acct_already_verified", "", requestedLanguage)})
		return
	}

	user.VerifiedAt = time.Now()
	user.VerifiedMethod = "Mail"

	if err := database.DB.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": translation.GetTranslation("update_verification_date_failed", "", requestedLanguage)})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "success.html", gin.H{"message": translation.GetTranslation("user_verified_successfully", "", requestedLanguage)})
}
