package handlers

import (
	"golang-gin-boilerplate/internal/controllers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func VoiceToTextHandler(c *gin.Context) {
	// Parse file from the form
	file, err := c.FormFile("audio_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve audio file"})
		return
	}

	// Save the uploaded file locally
	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save audio file"})
		return
	}
	defer os.Remove(filePath) // Clean up after processing

	// Process the file using the controller
	controller := controllers.VoiceToTextController{}
	text, err := controller.ConvertVoiceToText(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the recognized text
	c.JSON(http.StatusOK, gin.H{
		"recognized_text": text,
	})
}
