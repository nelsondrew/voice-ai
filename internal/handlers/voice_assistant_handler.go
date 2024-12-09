package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang-gin-boilerplate/internal/controllers"
	"golang-gin-boilerplate/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type VoiceAssistantHandler struct {
	voiceController  *controllers.VoiceToTextController
	chatController   *controllers.ChatGPTController
	context          *services.ConversationContext
	elevenLabsAPIKey string
}

func NewVoiceAssistantHandler(elevenLabsAPIKey string) *VoiceAssistantHandler {
	return &VoiceAssistantHandler{
		voiceController:  &controllers.VoiceToTextController{},
		chatController:   controllers.NewChatGPTController(),
		context:          services.NewConversationContext(openai.GPT3Dot5Turbo),
		elevenLabsAPIKey: elevenLabsAPIKey,
	}
}

func (h *VoiceAssistantHandler) VoiceAssistantHandler(c *gin.Context) {
	// Parse file from the form
	file, err := c.FormFile("audio_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve audio file",
		})
		return
	}

	// Save the uploaded file locally
	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save audio file",
		})
		return
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove temporary file: %v", err)
		}
	}()

	// Convert voice to text
	transcribedText, err := h.voiceController.ConvertVoiceToText(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to transcribe audio: " + err.Error(),
		})
		return
	}

	// Process transcribed text with ChatGPT
	assistantResponse, err := h.chatController.ProcessConversation(transcribedText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process conversation: " + err.Error(),
		})
		return
	}

	// Convert text to speech using Eleven Labs
	audioData, err := h.convertTextToSpeech(assistantResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert text to speech: " + err.Error(),
		})
		return
	}

	// Save the audio file temporarily
	audioFilePath := filepath.Join("/tmp", "assistant_response.mp3")
	if err := os.WriteFile(audioFilePath, audioData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save audio response",
		})
		return
	}
	defer func() {
		if err := os.Remove(audioFilePath); err != nil {
			log.Printf("Failed to remove temporary audio file: %v", err)
		}
	}()

	// Set the headers for audio response
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Content-Disposition", "inline; filename=assistant_response.mp3")

	// Return transcribed text, AI response, and audio file
	c.File(audioFilePath)
}

func (h *VoiceAssistantHandler) convertTextToSpeech(text string) ([]byte, error) {
	// Eleven Labs API endpoint (adjust the voice ID as needed)
	url := "https://api.elevenlabs.io/v1/text-to-speech/21m00Tcm4TlvDq8ikWAM"

	// Prepare the request body
	payload := map[string]interface{}{
		"text": text,
		"voice_settings": map[string]interface{}{
			"stability":        0.5,
			"similarity_boost": 0.5,
		},
	}

	// Create JSON payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", h.elevenLabsAPIKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Read the audio response
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return audioData, nil
}

// Existing reset and token handlers remain the same
func (h *VoiceAssistantHandler) ResetConversationHandler(c *gin.Context) {
	h.context.ResetContext()
	c.JSON(http.StatusOK, gin.H{
		"message": "Conversation context reset successfully",
	})
}

func (h *VoiceAssistantHandler) GetContextTokensHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_tokens": h.context.CalculateTotalTokens(),
	})
}
