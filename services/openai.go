package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService handles interactions with OpenAI API
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService creates a new OpenAIService instance
func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
	}
}

// ProcessText calls the OpenAI API to process a text part
func (s *OpenAIService) ProcessText(textPart string) (string, error) {
	resp, err := s.client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:     openai.GPT3TextDavinci003,
			Prompt:    textPart,
			MaxTokens: 150,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI: %v", err)
	}
	return strings.TrimSpace(resp.Choices[0].Text), nil
}

// SplitText splits text into parts of maximum length
func (s *OpenAIService) SplitText(text string, maxLength int) []string {
	var parts []string
	for len(text) > 0 {
		if len(text) <= maxLength {
			parts = append(parts, text)
			break
		}
		parts = append(parts, text[:maxLength])
		text = text[maxLength:]
	}
	return parts
} 