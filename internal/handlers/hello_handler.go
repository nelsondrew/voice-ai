package handlers

import (
	"golang-gin-boilerplate/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	// Use controller to get hello message
	message := controllers.GetHelloMessage()

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func CreateHelloHandler(c *gin.Context) {
	// Example of a POST handler
	var request struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := controllers.CreateHello(request.Name)

	c.JSON(http.StatusCreated, gin.H{
		"message": response,
	})
}
