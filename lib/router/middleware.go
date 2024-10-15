package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	model_token "github.com/roto17/zeus/lib/models/tokens"
	"github.com/roto17/zeus/lib/translation"
	"github.com/roto17/zeus/lib/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware checks for the JWT token in the Authorization header and verifies the role
func JWTAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("required_header", "", requested_language)})
			c.Abort()
			return
		}

		// Ensure the token has the correct Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token_format", "", requested_language)})
			c.Abort()
			return
		}

		// Extract the token from the Bearer <token> format
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(config.GetEnv("secretkey")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_token", "", requested_language)})
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
		userRole := claims["role"].(string)
		// Convert Unix timestamp back to time.Time for better handling
		expiration := time.Unix(int64(claims["exp"].(float64)), 0)

		// Check if the token has expired
		if time.Now().After(expiration) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("token_expired", "", requested_language)})
			c.Abort()
			return
		}

		var tokenRecord model_token.Token

		if err := database.DB.Where("token = ? and user_id = ?", tokenString, userId).Preload("User").First(&tokenRecord).Error; err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_or_expired_token", "", requested_language)})
			c.Abort()
			return
		}

		if tokenRecord.User.VerifiedAt.IsZero() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("verify_account", "", requested_language)})
			c.Abort()
			return
		}

		// Check if the role is allowed
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next() // Proceed to the next handler
				return
			}
		}

		// If the user's role is not allowed
		c.JSON(http.StatusForbidden, gin.H{"error": translation.GetTranslation("permission_denied", "", requested_language)})
		c.Abort()
	}
}

// Middleware to set a variable in the context
func SetHeaderVariableMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set a variable in the context
		c.Set("requested_language", utils.Coalesce(c.GetHeader("Accept-Language"), "en"))
		// Continue to the next handler
		c.Next()
	}
}
