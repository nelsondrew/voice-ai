package routes

import (
	"golang-gin-boilerplate/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	voiceAssistantHandler := handlers.NewVoiceAssistantHandler(os.Getenv("ELEVEN_LABS_API_KEY"))

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

	v1 := router.Group("/v1")
	{
		v1.POST("/voice-to-text", handlers.VoiceToTextHandler)
		v1.POST("/voice-assistant", voiceAssistantHandler.VoiceAssistantHandler)
	}

	return router
}
