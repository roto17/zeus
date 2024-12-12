package router

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	model_token "github.com/roto17/zeus/lib/models/tokens"
	"github.com/roto17/zeus/lib/sharedkeys"
	"github.com/roto17/zeus/lib/translation"
	"github.com/roto17/zeus/lib/utils"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

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

			result := database.DB.Where("token = ?", tokenString).Delete(&model_token.Token{})

			if result.RowsAffected == 0 || result.Error != nil {
				fmt.Printf("failed to delete token: %s", result.Error)
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": translation.GetTranslation("invalid_or_expired_token", "", requested_language)})
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

		c.Set(sharedkeys.UserKey, claims)

		// fmt.Printf("Claims being set: %#v\n", claims)

		// c.Set("user", claims)
		userId := claims["user_id"].(string)
		userRole := claims["role"].(string)

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

// Middleware to set a variable in the context
func SetEscapedID() gin.HandlerFunc {
	return func(c *gin.Context) {

		rawPath := c.Param("path") // Extract the path parameter

		// // Ensure rawPath is not empty
		// if len(rawPath) == 0 || rawPath[0] != '/' {
		// 	fmt.Printf("Invalid path")
		// 	// c.JSON(400, gin.H{"error": "Invalid path"})
		// 	return
		// }

		// Safely slice the rawPath
		if len(rawPath) != 0 && rawPath[0] == '/' {
			rawPath = rawPath[1:] // Remove the leading '/'
		}

		encodedPath := url.QueryEscape(rawPath)

		// if encodedPath == "" {
		// 	encodedPath = "A"
		// }

		// Set a variable in the context
		c.Set("escapedID", encodedPath)
		// Continue to the next handler
		c.Next()

	}
}

// rateLimiter returns a rate limiter for the given IP address.
func setLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Allow 5 requests per minute
		limiter = rate.NewLimiter(rate.Every(time.Minute/4), 4)
		visitors[ip] = limiter

		// Automatically delete the entry after 1 minute
		go func() {
			time.Sleep(time.Minute)
			mu.Lock()
			delete(visitors, ip)
			mu.Unlock()
		}()
	}

	return limiter
}

// RateLimitMiddleware is a middleware for rate limiting requests.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
		ip := c.ClientIP()

		limiter := setLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": translation.GetTranslation("too_many_request", "", requested_language),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
