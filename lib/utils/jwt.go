package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your_secret_key")

// GenerateJWT generates a JWT token for a specific user with their role
func GenerateJWT(username, role string) (string, error) {
	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
