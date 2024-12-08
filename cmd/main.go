package main

import (
	"golang-gin-boilerplate/internal/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func isGCPEnvironment() bool {
	// Check for Cloud Run specific environment variable
	return os.Getenv("K_SERVICE") != ""
}

func main() {

	isGCP := isGCPEnvironment()

	// If not running in GCP, load environment variables from .env file
	if !isGCP {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port := os.Getenv("CUSTOM_PORT")

	if port != "" {
		log.Println("Got the port:", port)
	} else {
		log.Println("PORT environment variable is not set")
	}
	// Initialize and start the server
	router := routes.SetupRouter()

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
