package inits

import (
    "context"
    "fmt"

    "google.golang.org/genai"
)

func NewGeminiClient(geminiApiKey string) (*genai.Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiApiKey,
        Backend: genai.BackendGeminiAPI,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create Gemini client: %w", err)
    }
	
    return client, nil
}