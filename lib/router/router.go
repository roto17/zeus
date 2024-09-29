package router

import (
	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/actions"
)

// var secretKey = []byte("your_secret_key")

func InitRouter() *gin.Engine {
	router := gin.Default()

	// Apply the middleware globally
	router.Use(SetHeaderVariableMiddleware())

	// Routes
	router.POST("/login", actions.Login)
	router.POST("/register", actions.Register)
	router.POST("/logout", actions.Logout)

	// Route for viewing a user by ID (Admin access only)
	router.GET("/view_user/:id", JWTAuthMiddleware("admin"), actions.ViewUser)

	return router
}
