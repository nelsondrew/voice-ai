package main

import (
	"golang-gin-boilerplate/internal/routes"
	"log"
)

func main() {
	// Initialize and start the server
	router := routes.SetupRouter()

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
