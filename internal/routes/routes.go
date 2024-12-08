package routes

import (
	"golang-gin-boilerplate/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Hello World routes
	helloGroup := router.Group("/hello")
	{
		helloGroup.GET("", handlers.HelloHandler)
		helloGroup.POST("/create", handlers.CreateHelloHandler)
	}

	// User routes (example)
	userGroup := router.Group("/users")
	{
		userGroup.GET("", handlers.GetUsersHandler)
		userGroup.POST("", handlers.CreateUserHandler)
	}

	return router
}
