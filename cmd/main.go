package main

import (
	"fmt"
	"golang-gin-boilerplate/internal/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Loaded env")
	// Initialize and start the server
	router := routes.SetupRouter()

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
