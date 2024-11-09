package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roto17/zeus/lib/config"
	model_user "github.com/roto17/zeus/lib/models/users"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GenerateToken generates a JWT for the authenticated user
func GenerateToken(user model_user.User) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Hour * 24) // Token expires in 24 hours

	// Create token claims including the role and expiration time
	claims := jwt.MapClaims{
		"user_id":  fmt.Sprintf("%v", user.ID),
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
		// "verified": !user.VerifiedAt.IsZero(),
	}

	// fmt.Printf("-------------------------\n")
	// fmt.Printf("%v", user.VerifiedAt.Second())
	// fmt.Printf("-------------------------\n")

	// Create the token with signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// fmt.Printf("%s", []byte(config.GetEnv("secretkey")))

	// Sign and return the token
	tokenString, err := token.SignedString([]byte(config.GetEnv("secretkey")))
	if err != nil {
		return "", err
	}

	// // Store the token in the Token table
	// tokenEntry := model_token.Token{
	// 	Token:     tokenString,
	// 	ExpiresAt: expirationTime,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	// if err := database.DB.Create(&tokenEntry).Error; err != nil {
	// 	return "", expirationTime, err
	// }

	return tokenString, nil
}

// Hash and compare the password using bcrypt
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetRoleFromToken extracts the role from a JWT token in the Authorization header.
func GetRoleFromToken(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("secretkey")), nil
	})
	if err != nil || !token.Valid {
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if role, exists := claims["role"]; exists {
			return role.(string)
		}
	}
	return ""
}

// EncryptID encrypts an integer ID using AES encryption
func EncryptID(id uint) string {
	// Convert the ID to a byte slice
	idBytes := []byte(fmt.Sprintf("%d", id))

	// Create a new AES cipher with the provided key
	block, err := aes.NewCipher([]byte(config.GetEnv("encryption_key")))
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Generate a new GCM (Galois/Counter Mode) cipher based on AES
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Printf("%v", err)
	}

	// Encrypt the ID bytes with AES-GCM
	ciphertext := aesGCM.Seal(nonce, nonce, idBytes, nil)

	// Encode the ciphertext to base64 for easier storage or transmission
	return url.QueryEscape(base64.StdEncoding.EncodeToString(ciphertext))
}

// DecryptID decrypts the encrypted string to retrieve the original integer ID
func DecryptID(encryptedID string) uint {

	// URL decode the string
	decodedStr, err := url.QueryUnescape(encryptedID)
	if err != nil {
		// Handle error if decoding fails
		fmt.Println("Error decoding:", err)
	}

	// Decode the base64-encoded string
	ciphertext, err := base64.StdEncoding.DecodeString(decodedStr)
	if err != nil {
		fmt.Printf("test %v", err)
	}

	// Create a new AES cipher with the provided key
	block, err := aes.NewCipher([]byte(config.GetEnv("encryption_key")))
	if err != nil {
		fmt.Printf("error key %v", err)
	}

	// Generate a new GCM cipher based on AES
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		// return 0, err
		fmt.Printf("test %v", err)
	}

	// Split the nonce and the actual ciphertext
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	idBytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// return 0, err
		fmt.Printf("tes %v", err)
	}

	// Convert the decrypted bytes back to an integer
	var id uint
	fmt.Sscanf(string(idBytes), "%d", &id)
	return id
}
