package services

import (
	"sync"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ConversationMessageType defines the type of message in the conversation
type ConversationMessageType string

const (
	MessageTypeSystem    ConversationMessageType = "system"
	MessageTypeUser      ConversationMessageType = "user"
	MessageTypeAssistant ConversationMessageType = "assistant"
)

// ConversationContext manages the entire conversation state
type ConversationContext struct {
	mu                 sync.RWMutex
	Messages           []openai.ChatCompletionMessage
	MaxTokens          int
	CurrentModel       string
	LanguagePreference language.Tag
}

// NewConversationContext creates a new conversation context
func NewConversationContext(model string) *ConversationContext {
	return &ConversationContext{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful AI assistant. Maintain context of our ongoing conversation.",
			},
		},
		MaxTokens:          4096,
		CurrentModel:       model,
		LanguagePreference: language.English,
	}
}

// AddMessage adds a new message to the conversation context
func (cc *ConversationContext) AddMessage(
	messageType ConversationMessageType,
	content string,
) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	newMessage := openai.ChatCompletionMessage{
		Role:    string(messageType),
		Content: content,
	}
	cc.Messages = append(cc.Messages, newMessage)

	// Automatically trim context if needed
	cc.trimContextUnlocked() // Use an unlocked version
}

// trimContextUnlocked is an internal method that assumes the mutex is already held
func (cc *ConversationContext) trimContextUnlocked() {
	// If total tokens are within acceptable limit, do nothing
	if cc.calculateTotalTokens() <= cc.MaxTokens {
		return
	}

	// Always keep the system message
	systemMessage := cc.Messages[0]

	// Create a new slice to store trimmed messages
	var trimmedMessages []openai.ChatCompletionMessage
	trimmedMessages = append(trimmedMessages, systemMessage)

	// Work backwards through conversation history
	for i := len(cc.Messages) - 1; i > 0; i-- {
		currentMessage := cc.Messages[i]

		// Add message to trimmed messages
		trimmedMessages = append([]openai.ChatCompletionMessage{currentMessage}, trimmedMessages...)

		// Recalculate total tokens after adding this message
		if cc.calculateTokensForMessageList(trimmedMessages) <= cc.MaxTokens {
			break
		}
	}

	// Update the messages
	cc.Messages = trimmedMessages
}

// TrimContext is a public method that can be called externally
func (cc *ConversationContext) TrimContext() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.trimContextUnlocked()
}

// CalculateTokensForMessage estimates tokens for a single message
func (cc *ConversationContext) calculateTokensForMessage(
	message openai.ChatCompletionMessage,
) int {
	// Rough estimation: ~4 characters per token
	baseTokens := len(message.Content) / 4

	// Add extra tokens based on role
	switch message.Role {
	case "system":
		baseTokens += 10 // More weight for system messages
	case "user":
		baseTokens += 5
	case "assistant":
		baseTokens += 5
	}

	return baseTokens
}

// CalculateTokensForMessageList calculates tokens for multiple messages
func (cc *ConversationContext) calculateTokensForMessageList(
	messages []openai.ChatCompletionMessage,
) int {
	totalTokens := 0
	for _, msg := range messages {
		totalTokens += cc.calculateTokensForMessage(msg)
	}
	return totalTokens
}

// CalculateTotalTokens calculates tokens in current context
func (cc *ConversationContext) calculateTotalTokens() int {
	return cc.calculateTokensForMessageList(cc.Messages)
}

// SetLanguagePreference allows setting conversation language
func (cc *ConversationContext) SetLanguagePreference(lang language.Tag) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.LanguagePreference = lang
}

// GetLocalizedMessage provides localized system messages
func (cc *ConversationContext) GetLocalizedMessage(key string) string {
	p := message.NewPrinter(cc.LanguagePreference)

	// Example localization map (you'd expand this)
	messages := map[string]string{
		"welcome": p.Sprintf("Welcome! How can I assist you today?"),
		"help":    p.Sprintf("I'm here to help. What do you need?"),
	}

	return messages[key]
}

// ResetContext completely resets the conversation
func (cc *ConversationContext) ResetContext() {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.Messages = []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are a helpful AI assistant. Maintain context of our ongoing conversation.",
		},
	}
}

func (cc *ConversationContext) CalculateTotalTokens() int {
	if cc == nil {
		return 0
	}

	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if len(cc.Messages) == 0 {
		return 0
	}

	totalTokens := 0
	for _, msg := range cc.Messages {
		totalTokens += cc.calculateTokensForMessage(msg)
	}
	return totalTokens
}
