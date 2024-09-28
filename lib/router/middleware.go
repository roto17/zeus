package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/models"
	"github.com/roto17/zeus/lib/translation"
	"github.com/roto17/zeus/lib/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// var secretKey = []byte("your_secret_key")

// JWTAuthMiddleware checks for the JWT token in the Authorization header and verifies the role
func JWTAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := utils.Coalesce(c.GetHeader("Accept-Language"), "en")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("required_header", "", lang)})
			c.Abort()
			return
		}

		// Ensure the token has the correct Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token_format", "", lang)})
			c.Abort()
			return
		}

		// Extract the token from the Bearer <token> format
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Find the token in the database
		var tokenRecord models.Token

		if err := database.DB.Where("token = ?", tokenString).First(&tokenRecord).Error; err != nil || tokenRecord.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_or_expired_token", "", lang)})
			c.Abort()
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token", "", lang)})
			c.Abort()
			return
		}

		// Extract the claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token_claims", "", lang)})
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
		c.JSON(http.StatusForbidden, gin.H{"error": translation.GetTranslation("permission_denied", "", lang)})
		c.Abort()
	}
}
