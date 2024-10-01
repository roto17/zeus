package actions

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	model_token "github.com/roto17/zeus/lib/models/tokens"
	model_user "github.com/roto17/zeus/lib/models/users"
	"github.com/roto17/zeus/lib/translation" // Assuming translation package handles translations
	"github.com/roto17/zeus/lib/utils"
)

func CreateUser(user *model_user.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

func GetUser(id int) (model_user.User, error) {
	var user model_user.User
	result := database.DB.First(&user, id)
	return user, result.Error
}

func UpdateUser(user *model_user.User) error {
	result := database.DB.Save(&user)
	return result.Error
}

func DeleteUser(id int) error {
	result := database.DB.Delete(&model_user.User{}, id)
	return result.Error
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

// Login handles login and JWT generation
func Login(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB

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
	if err := db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("unauthorized", "", requested_language)})
		return
	}

	// Check if the password matches
	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("unauthorized", "", requested_language)})
		return
	}

	// Generate JWT token
	token, expiration, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("token_generation_failed", "", requested_language)})
		return
	}

	// Save the token and expiration in the database
	newToken := model_token.Token{
		Token:     token,
		ExpiresAt: expiration,
	}
	if err := db.Create(&newToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("token_save_failed", "", requested_language)})
		return
	}

	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

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

	// // Create a new user record
	// newUser := model_user.User{
	// 	Username: user.Username,
	// 	Password: hashedPassword,
	// 	Role:     user.Role,
	// }

	newUser := model_user.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Password:  hashedPassword,
		Role:      user.Role,
		// VerifiedAt: user.VerifiedAt,
		MiddleName: user.MiddleName,
	}

	// Save the user in the database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("registration_failed", "", requested_language)})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("registration_successful", "", requested_language)})
}

// Logout handles logout by invalidating the JWT token
func Logout(c *gin.Context) {
	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
	db := database.DB

	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("token_not_provided", "", requested_language)})
		return
	}

	// Strip "Bearer " if it's included in the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Find the token in the database
	var token model_token.Token
	if err := db.Where("token = ?", tokenString).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token", "", requested_language)})
		return
	}

	// Set the expiration to the current time to invalidate the token
	token.ExpiresAt = time.Now()
	if err := db.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": translation.GetTranslation("token_invalidation_failed", "", requested_language)})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": translation.GetTranslation("logout_successful", "", requested_language)})
}
