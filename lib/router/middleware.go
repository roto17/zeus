package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// var secretKey = []byte("your_secret_key")

// JWTAuthMiddleware checks for the JWT token in the Authorization header and verifies the role
func JWTAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Ensure the token has the correct Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extract the token from the Bearer <token> format
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))

		// Find the token in the database
		var tokenRecord models.Token
		fmt.Printf("%s", tokenString)
		// fmt.Printf("tzzzzzzzzzzzz")
		if err := database.DB.Where("token = ?", tokenString).First(&tokenRecord).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if the token is expired
		if tokenRecord.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract the claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Check if the role is allowed
		userRole := claims["role"].(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next() // Proceed to the next handler
				return
			}
		}

		// If the user's role is not allowed
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
		c.Abort()
	}
}
