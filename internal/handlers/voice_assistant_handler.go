package handlers

import (
	"log"
	"net/http"
	"os"

	"golang-gin-boilerplate/internal/controllers"
	"golang-gin-boilerplate/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type VoiceAssistantHandler struct {
	voiceController *controllers.VoiceToTextController
	chatController  *controllers.ChatGPTController
	context         *services.ConversationContext
}

func NewVoiceAssistantHandler() *VoiceAssistantHandler {
	return &VoiceAssistantHandler{
		voiceController: &controllers.VoiceToTextController{},
		chatController:  controllers.NewChatGPTController(),
		context:         services.NewConversationContext(openai.GPT3Dot5Turbo),
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

	// Return both transcribed text and AI response
	c.JSON(http.StatusOK, gin.H{
		"transcribed_text":     transcribedText,
		"assistant_response":   assistantResponse,
		"total_context_tokens": h.context.CalculateTotalTokens(),
	})
}

// Additional handler methods for conversation management
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
