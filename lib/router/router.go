package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/actions"
	"github.com/roto17/zeus/lib/notifications"
	"github.com/roto17/zeus/lib/translation"
	"github.com/roto17/zeus/lib/utils"
)

// InitRouter initializes the Gin router and starts the necessary components
func InitRouter(ctx context.Context) *gin.Engine {
	router := gin.Default()

	// Start 5 worker goroutines for handling notifications, using the provided context
	notifications.StartWorkers(5, ctx)

	// Load HTML templates for error pages
	router.LoadHTMLGlob("lib/views/*/*")

	// Apply middleware
	router.Use(SetHeaderVariableMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust to your needs
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Forwarded-For"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Group routes under the /api prefix
	api := router.Group("/api")
	{
		// User-related routes
		api.POST("/register", actions.Register)
		api.POST("/login", RateLimitMiddleware(), actions.Login)
		api.POST("/logout", actions.Logout)
		api.POST("/logout-all", actions.LogoutAll)

		// WebSocket route for notifications
		api.GET("/notifications", notifications.WSHandler)

		// Route for viewing a user by ID (Admin access only)
		api.GET("/view_user/:id", JWTAuthMiddleware("admin"), actions.ViewUser)

		api.POST("/product_categories", JWTAuthMiddleware("admin"), actions.AddProductCategory)
		api.PATCH("/product_categories", JWTAuthMiddleware("admin"), actions.UpdateProductCategory)
		api.DELETE("/product_categories/:id", JWTAuthMiddleware("admin"), actions.DeleteProductCategory)

		api.POST("/products", JWTAuthMiddleware("admin"), actions.AddProduct)

		// Other routes
		router.GET("/verify-email", actions.VerifyByMail)
	}

	// Handle undefined routes (still under the /api prefix)
	router.NoRoute(func(c *gin.Context) {
		requestedLanguage := utils.GetHeaderVarToString(c.Get("requested_language"))
		errorMessage := translation.GetTranslation("not_found", "", requestedLanguage)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
	})

	go notifications.HandleMessages() // Start the notification message handler

	return router
}
