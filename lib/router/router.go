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
	router.Use(SetEscapedID())

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
		api.POST("/users", actions.Register)
		api.PATCH("/users", JWTAuthMiddleware("admin", "super_admin"), actions.UpdateUser)
		api.GET("/users/*path", JWTAuthMiddleware("admin", "super_admin"), actions.ViewUser)
		api.GET("/users", JWTAuthMiddleware("admin"), actions.AllUsers)
		api.DELETE("/users/*path", JWTAuthMiddleware("admin"), actions.DeleteUser)

		api.POST("/companies", actions.AddCompany)
		api.PATCH("/companies", actions.UpdateCompany)
		api.GET("/companies/*path", actions.ViewCompany)
		api.GET("/companies", actions.AllCompanies)
		api.DELETE("/companies/*path", actions.DeleteCompany)

		api.POST("/login", RateLimitMiddleware(), actions.Login)
		router.GET("/verify-email", actions.VerifyByMail)
		api.POST("/logout", JWTAuthMiddleware("admin", "super_admin"), actions.Logout)
		api.POST("/logout-all", JWTAuthMiddleware("admin", "super_admin"), actions.LogoutAll)

		// api.GET("/verify_sms", JWTAuthMiddleware("admin", "super_admin"), actions.VerifyBySMS)

		api.POST("/product_categories", JWTAuthMiddleware("admin"), actions.AddProductCategory)
		api.PATCH("/product_categories", JWTAuthMiddleware("admin"), actions.UpdateProductCategory)
		api.GET("/product_categories/*path", JWTAuthMiddleware("admin"), actions.ViewProductCategory)
		api.GET("/product_categories", JWTAuthMiddleware("admin"), actions.AllProductCategories)
		api.DELETE("/product_categories/*path", JWTAuthMiddleware("admin"), actions.DeleteProductCategory)

		api.POST("/products", JWTAuthMiddleware("admin"), actions.AddProduct)
		api.PATCH("/products", JWTAuthMiddleware("admin"), actions.UpdateProduct)
		// api.PATCH("/products", JWTAuthMiddleware("admin"), actions.UpdateProductTest)
		api.GET("/products/*path", JWTAuthMiddleware("admin"), actions.ViewProduct)
		api.GET("/products", JWTAuthMiddleware("admin"), actions.AllProducts)
		api.DELETE("/products/*path", JWTAuthMiddleware("admin"), actions.DeleteProduct)

		api.POST("/orders", JWTAuthMiddleware("admin"), actions.AddOrder)
		api.POST("/orders_add_product", JWTAuthMiddleware("admin"), actions.AddProductToStock)

		// api.GET("/products", JWTAuthMiddleware("admin"), actions.SaveQR)

		// WebSocket route for notifications
		api.GET("/notifications", notifications.WSHandler)
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
