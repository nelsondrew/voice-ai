package controllers

import (
	"fmt"
	"golang-gin-boilerplate/internal/interfaces"
	"golang-gin-boilerplate/internal/models"
)

type HelloController struct {
	HelloInterface interfaces.HelloInterface
}

func GetHelloMessage() string {
	return "Hello, World from Gin Backend!"
}

func CreateHello(name string) string {
	// Create a hello model
	helloModel := models.HelloModel{
		Name: name,
	}

	return fmt.Sprintf("Created hello for %s", helloModel.Name)
}
