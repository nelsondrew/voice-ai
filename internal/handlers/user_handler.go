package handlers

import (
	"golang-gin-boilerplate/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsersHandler(c *gin.Context) {
	users := controllers.GetUsers()
	c.JSON(http.StatusOK, users)
}

func CreateUserHandler(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdUser := controllers.CreateUser(user.Username, user.Email)
	c.JSON(http.StatusCreated, createdUser)
}
