package inits

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Secrets struct {
	DatabaseURL      string
	GeminiAPIKey   	 string
	S3Region		 string
}

func InitSecrets() (*Secrets, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, trying with inline environment variables anyway")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		return nil, fmt.Errorf("S3_REGION environment variable is required")
	}

	return &Secrets{
		DatabaseURL: databaseURL,
		GeminiAPIKey: geminiAPIKey,
		S3Region: s3Region,
	}, nil
}