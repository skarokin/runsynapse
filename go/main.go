package main

import (
	"log"
	"net/http"
	"os"

	"github.com/skarokin/runsynapse/go/inits"
	"github.com/skarokin/runsynapse/go/handlers"
)

func main() {
	secrets, err := inits.InitSecrets()
	if err != nil {
		log.Fatalf("Failed to initialize secrets: %v", err)
	}

	supabaseClient, err := inits.NewSupabaseClient(secrets.DatabaseURL) 
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}
	defer supabaseClient.Close()

	geminiClient, err := inits.NewGeminiClient(secrets.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	handler := handlers.NewHandler(supabaseClient, geminiClient)

	handler.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

    log.Printf("Server running on port %s", port)
    log.Printf("Health check: http://localhost:%s/health", port)

    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}