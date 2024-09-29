package actions

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/utils"
)

func CreateUser(user *models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

func GetUser(id int) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, id)
	return user, result.Error
}

func UpdateUser(user *models.User) error {
	result := database.DB.Save(&user)
	return result.Error
}

func DeleteUser(id int) error {
	result := database.DB.Delete(&models.User{}, id)
	return result.Error
}

// ViewUser handler
func ViewUser(c *gin.Context) {
	// Get the user ID from URL param
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Use the GetUser function to fetch the user by ID
	user, err := GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user
	c.JSON(http.StatusOK, user)
}

// Login handles user login and JWT generation
func Login(c *gin.Context) {
	db := database.DB

	// Extract credentials from request
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user by username
	var user models.User
	if err := db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check if the password matches
	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, expiration, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Save the token and expiration in the database
	newToken := models.Token{
		Token:     token,
		ExpiresAt: expiration,
	}
	if err := db.Create(&newToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save token"})
		return
	}

	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Register handles user registration
func Register(c *gin.Context) {

	requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

	db := database.DB

	var user models.User

	// Bind the incoming JSON to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate and get translated error messages
	validationErrors := utils.FieldValidationAll(user, requested_language)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create a new user record
	newUser := models.User{
		Username: user.Username,
		Password: hashedPassword,
		Role:     user.Role,
	}

	// Save the user in the database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Logout handles user logout by invalidating the JWT token
func Logout(c *gin.Context) {
	db := database.DB

	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
		return
	}

	// Strip "Bearer " if it's included in the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Find the token in the database
	var token models.Token
	if err := db.Where("token = ?", tokenString).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Set the expiration to the current time to invalidate the token
	token.ExpiresAt = time.Now()
	if err := db.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not invalidate token"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
