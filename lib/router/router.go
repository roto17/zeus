package router

import (
	"github.com/gin-gonic/gin"
	"github.com/roto17/zeus/lib/actions"
)

// var secretKey = []byte("your_secret_key")

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Routes
	r.POST("/login", actions.Login)
	r.POST("/register", actions.Register)
	r.POST("/logout", actions.Logout)

	// Route for viewing a user by ID (Admin access only)
	r.GET("/view_user/:id", JWTAuthMiddleware("admin"), actions.ViewUser)

	return r
}
