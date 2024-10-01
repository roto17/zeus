package router

import (
	"net/http"

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

	// Group routes under the /api prefix
	api := router.Group("/api")
	// Define routes within the /api prefix
	api.POST("/register", actions.Register)
	api.POST("/login", actions.Login)
	api.POST("/logout", actions.Logout)

	// Route for viewing a user by ID (Admin access only)
	api.GET("/view_user/:id", JWTAuthMiddleware("admin"), actions.ViewUser)

	// Handle undefined routes (still under the /api prefix)
	router.NoRoute(func(c *gin.Context) {
		requested_language := utils.GetHeaderVarToString(c.Get("requested_language"))
		errorMessage := translation.GetTranslation("not_found", "", requested_language)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
	})

	return router

}
