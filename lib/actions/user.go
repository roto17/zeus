package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	model_token "github.com/roto17/zeus/lib/models/tokens"
	model_user "github.com/roto17/zeus/lib/models/users"
	"github.com/roto17/zeus/lib/notifications"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
	"gopkg.in/gomail.v2"
)

// Register handles registration
func Register(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB
	var new_user model_user.CreateUserInput

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&new_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_input", "", requested_language)})
		return
	}

	user := model_user.User{
		Email:      new_user.Email,
		FirstName:  new_user.FirstName,
		LastName:   new_user.LastName,
		Username:   new_user.Username,
		Password:   new_user.Password,
		Role:       new_user.Role,
		VerifiedAt: new_user.VerifiedAt,
		MiddleName: new_user.MiddleName,
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(user, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("password_hashing_failed", "", requested_language)})
		return
	}

	newUser := model_user.User{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Password:   hashedPassword,
		Role:       user.Role,
		MiddleName: user.MiddleName,
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

	if err := SendVerificationEmail(user.Email, token, config.GetEnv("appBaseURL"), config.GetEnv("smtpUser"), config.GetEnv("smtpPass"), config.GetEnv("smtpHost"), utils.StringToInt((config.GetEnv("smtpPort")))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("failed_to_send_verification_email", "", requested_language)})
		return
	}
	// Return success message
	// notifications.SendNotification("A new user has been added: " + user.Username)

	// After successfully registering a user
	fromRole := "user"
	toRole := "user"                                               // Define the user's role
	notifications.RegisterUser(fromRole, toRole, "added success!") // Call the RegisterUser function to send a notification

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

// ViewUser handler
func ViewUser(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	// Get the user ID from URL param
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("invalid_id", "", requested_language)})
		return
	}

	// Use the GetUser function to fetch the user by ID
	user, err := GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
		return
	}

	// Return the user
	c.JSON(http.StatusOK, user)
}

// func CreateUser(user *model_user.User) error {
// 	result := database.DB.Create(&user)
// 	return result.Error
// }

func GetUser(id int) (model_user.User, error) {
	var user model_user.User
	result := database.DB.First(&user, id)
	return user, result.Error
}

// func UpdateUser(user *model_user.User) error {
// 	result := database.DB.Save(&user)
// 	return result.Error
// }

// func DeleteUser(id int) error {
// 	result := database.DB.Delete(&model_user.User{}, id)
// 	return result.Error
// }

func SendVerificationEmail(userEmail, token, appBaseURL, smtpUser, smtpPass, smtpHost string, smtpPort int) error {
	// Create the verification URL
	verificationURL := fmt.Sprintf("%s/verify-email?signature=%s", appBaseURL, token)

	// Email content
	subject := "Email Verification"
	// body := fmt.Sprintf("Please click the following link to verify your email: %s", verificationURL)

	// HTML version of the email body
	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<p>Please click the following link to verify your email:</p>
				<a href="%s">Verify Email</a>
			</body>
		</html>`, verificationURL)

	// Set up the email message
	message := gomail.NewMessage()
	message.SetHeader("From", smtpUser)
	message.SetHeader("To", userEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", htmlBody) // plain text
	// message.AddAlternative("text/html", htmlBody) // HTML version

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

// // Register handles registration
// func VerifyByMail(c *gin.Context) {
// 	tokenString := c.Query("signature")
// 	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

// 	if tokenString == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": translation.GetTranslation("signature_required", "", requested_language)})
// 		return
// 	}

// 	// Parse the token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, http.ErrAbortHandler
// 		}
// 		return []byte(config.GetEnv("secretkey")), nil
// 	})

// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_or_expired_token", "", requested_language)})
// 		c.Abort()
// 		return
// 	}

// 	// Extract the claims
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token_claims", "", requested_language)})
// 		c.Abort()
// 		return
// 	}

// 	user_id := utils.StringToInt(claims["user_id"].(string))

// 	var user model_user.User
// 	result := database.DB.First(&user, user_id)

// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": translation.GetTranslation("not_found", "", requested_language)})
// 		c.Abort()
// 		return
// 	}

// 	if !user.VerifiedAt.IsZero() {
// 		c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("acct_already_verified", "", requested_language)})
// 		return

// 	}

// 	user.VerifiedAt = time.Now()

// 	if err := database.DB.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("update_verification_date_failed", "", requested_language)})
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("user_verified_successfully", "", requested_language)})

// }

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

	if err := database.DB.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": translation.GetTranslation("update_verification_date_failed", "", requestedLanguage)})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "success.html", gin.H{"message": translation.GetTranslation("user_verified_successfully", "", requestedLanguage)})
}
