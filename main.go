// File: main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/config" // Replace with your actual module path
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/router"
	"github.com/roto17/zeus/lib/utils"
)

// Secret key for signing the token (keep it safe!)
var jwtKey = []byte("your_secret_key")

// Define your JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Function to generate the token
func GenerateJWT(username string) (string, error) {
	// Set expiration time (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the username and expiration time
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Handler function that works with gin.HandlerFunc
func GetUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	user, err := actions.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func main() {

	config.LoadConfig()

	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en" // Default to English
	}

	// Initialize the database
	database.InitDB()

	// migrations.MigrateUser()

	// Example of generating a token for a user
	// token, err := GenerateJWT("exampleuser")
	// if err != nil {
	// 	log.Fatal("Error generating token:", err)
	// }

	// log.Println("Generated JWT Token:", token)

	// user := models.User{Name: "John Doe", Desc: "Desc", Jam: "OKOK", Email: "test@test.com"}

	// // Validate and get translated error messages
	// validationerrors := utils.FieldValidationAll(user, "en")

	// if validationerrors != nil {
	// 	for _, err := range validationerrors {
	// 		fmt.Printf("%s: %s\n", err.Field, err.Message)
	// 	}
	// } else {
	// 	fmt.Printf("All is good")
	// }
	token, err := utils.GenerateJWT("roto", "admin")
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
	} else {
		fmt.Printf("Generated JWT Token: %s\n", token)
	}
	router.InitRouter().Run(":8080")

}
