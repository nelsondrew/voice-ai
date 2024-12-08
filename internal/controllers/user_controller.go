package controllers

import (
	"golang-gin-boilerplate/internal/models"
)

func GetUsers() []models.UserModel {
	return []models.UserModel{
		{
			ID:       1,
			Username: "john_doe",
			Email:    "john@example.com",
		},
		{
			ID:       2,
			Username: "jane_smith",
			Email:    "jane@example.com",
		},
	}
}

func CreateUser(username, email string) models.UserModel {
	return models.UserModel{
		ID:       3, // In real app, this would be generated
		Username: username,
		Email:    email,
	}
}
