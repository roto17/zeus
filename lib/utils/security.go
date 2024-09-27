package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/models"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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

// Hash and compare the password using bcrypt
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
