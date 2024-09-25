package router

import (
	"net/http"
	"strconv"

	"github.com/roto17/zeus/lib/actions" // Adjust based on your project structure

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	usersGroup := router.Group("/users")
	{
		// usersGroup.Use(JWTAuthMiddleware())
		// usersGroup.GET("/:id", GetUser)
		usersGroup.GET("/:id", JWTAuthMiddleware("admin"), GetUser)
	}

	return router
}

func GetUser(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, fetchErr := actions.GetUser(id)
	if fetchErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
