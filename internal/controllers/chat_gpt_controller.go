package controllers

import (
	"context"
	"fmt"
	"golang-gin-boilerplate/internal/services"
	"os"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTController struct {
	client  *openai.Client
	context *services.ConversationContext
}

func NewChatGPTController() *ChatGPTController {
	apiKey := os.Getenv("OPEN_API_KEY")
	if apiKey == "" {
		// Logging or error handling for missing API key
		panic("OPENAI_API_KEY environment variable is not set")
	}

	return &ChatGPTController{
		client:  openai.NewClient(apiKey),
		context: services.NewConversationContext(openai.GPT3Dot5Turbo),
	}
}

func (c *ChatGPTController) ProcessConversation(userInput string) (string, error) {
	// Add user message
	c.context.AddMessage(services.MessageTypeUser, userInput)

	totalTokens := 0
	for _, message := range c.context.Messages {
		totalTokens += len(message.Content) // Approximate token count, a better approach would involve using a tokenizer
	}

	fmt.Printf("Total tokens: %d\n", totalTokens)

	if len(c.context.Messages) > 0 {
		c.context.Messages = c.context.Messages[len(c.context.Messages)-1:]
	}

	// Prepare request
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		Messages:  c.context.Messages,
		MaxTokens: 1500,
	}

	// Get response from OpenAI
	resp, err := c.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %v", err)
	}

	// Extract and add assistant response
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	responseText := resp.Choices[0].Message.Content
	c.context.AddMessage(services.MessageTypeAssistant, responseText)

	return responseText, nil
}

// Optional: Method to reset conversation context
func (c *ChatGPTController) ResetConversation() {
	c.context.ResetContext()
}

// Optional: Method to get current context tokens
func (c *ChatGPTController) GetCurrentTokenCount() int {
	return c.context.CalculateTotalTokens()
}
