package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/translation"
	"github.com/roto17/zeus/lib/utils"
)

// var secretKey = []byte("your_secret_key")

func InitRouter() *gin.Engine {
	router := gin.Default()

	// Apply the middleware globally
	router.Use(SetHeaderVariableMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Add allowed origins here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Group routes under the /api prefix
	api := router.Group("/api")
	// Define routes within the /api prefix
	api.POST("/register", actions.Register)
	api.POST("/login", RateLimitMiddleware(), actions.Login)
	api.POST("/logout", actions.Logout)
	api.POST("/logout-all", actions.LogoutAll)

	// Route for viewing a user by ID (Admin access only)
	api.GET("/view_user/:id", JWTAuthMiddleware("admin"), actions.ViewUser)
	router.GET("verify-email", actions.VerifyByMail)

	// Handle undefined routes (still under the /api prefix)
	router.NoRoute(func(c *gin.Context) {
		requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
		errorMessage := translation.GetTranslation("not_found", "", requested_language)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
	})

	return router

}
