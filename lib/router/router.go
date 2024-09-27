package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/models"
	"golang.org/x/crypto/bcrypt"
)

// var secretKey = []byte("your_secret_key")

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Routes
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.POST("/logout", Logout)

	// Route for viewing a user by ID (Admin access only)
	r.GET("/view_user/:id", JWTAuthMiddleware("admin"), ViewUser)

	return r
}

// Hash and compare the password using bcrypt
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT for the authenticated user
func GenerateToken(user models.User) (string, time.Time, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Hour * 72) // Token expires in 72 hours

	// Create token claims including the role and expiration time
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	}

	// Create the token with signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Printf("%s", []byte(config.GetEnv("secretkey")))

	// Sign and return the token
	tokenString, err := token.SignedString([]byte(config.GetEnv("secretkey")))
	if err != nil {
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
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
	if !checkPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, expiration, err := GenerateToken(user)
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

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Register handles user registration
func Register(c *gin.Context) {
	db := database.DB

	// Extract user details from request
	var registerData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the user's password
	hashedPassword, err := HashPassword(registerData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create a new user record
	newUser := models.User{
		Username: registerData.Username,
		Password: hashedPassword,
		Role:     registerData.Role,
	}

	// Save the user in the database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
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
	user, err := actions.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user
	c.JSON(http.StatusOK, user)
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
