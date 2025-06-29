package utils

import (
    "context"
    "fmt"
    "log"
    "strings"
    "time"

    "google.golang.org/genai"
)

func getEmbedding(ctx context.Context, client *genai.Client, text string, taskType string) (string, error) {
	// generic embedding generator that can be used for both thoughts and queries
	// takes a task type to differentiate between retrieval and other tasks
	start := time.Now()
	log.Printf("[EMBEDDING] Starting embedding generation for text (length: %d chars)", len(text))
	
	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}

	config := &genai.EmbedContentConfig{
		TaskType: taskType,
	}

	log.Printf("[EMBEDDING] Calling Gemini API with model: gemini-embedding-exp-03-07")
	
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-exp-03-07",
		contents,
		config,
	)
	
	apiDuration := time.Since(start)
	
	if err != nil {
		log.Printf("[EMBEDDING] ERROR: API call failed after %v: %v", apiDuration, err)
		return "", fmt.Errorf("failed to get embedding: %w", err)
	}

	log.Printf("[EMBEDDING] API call successful in %v", apiDuration)

	if len(result.Embeddings) == 0 {
		log.Printf("[EMBEDDING] ERROR: No embeddings returned in response")
		return "", fmt.Errorf("no embeddings returned")
	}
	
	if len(result.Embeddings[0].Values) == 0 {
		log.Printf("[EMBEDDING] ERROR: Empty embedding values returned")
		return "", fmt.Errorf("no embedding values returned")
	}

	embedding := result.Embeddings[0].Values
	
	// convert to PostgreSQL vector format: [1.0,2.0,3.0]
	vectorStr := make([]string, len(embedding))
	for i, val := range embedding {
		vectorStr[i] = fmt.Sprintf("%.6f", val) // limit precision a bit to reduce size
	}

	vectorString := "[" + strings.Join(vectorStr, ",") + "]"
	
	totalDuration := time.Since(start)
	log.Printf("[EMBEDDING] Total embedding generation completed in %v", totalDuration)
	
	return vectorString, nil
}

// in this case pretty important to track embedding API call duration
func GetThoughtEmbedding(ctx context.Context, client *genai.Client, text string) (string, error) {
    log.Printf("[EMBEDDING] Generating embedding for thought: %s", text)

	return getEmbedding(ctx, client, text, "DOCUMENT_RETRIEVAL")
}

func GetQueryEmbedding(ctx context.Context, client *genai.Client, query string) (string, error) {
	log.Printf("[EMBEDDING] Generating embedding for query: %s", query)

	return getEmbedding(ctx, client, query, "QUERY_RETRIEVAL")
}