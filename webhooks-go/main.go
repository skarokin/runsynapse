package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"

    "github.com/skarokin/runsynapse/webhooks/supabase"
    "github.com/skarokin/runsynapse/webhooks/events"
)

type WebhookHandler struct {
	supabaseClient *supabase.SupabaseClient
	webhookSecret string
}

// smee --url [SMEE_URL] --path /webhook --port 8080
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("WEBHOOK_SECRET environment variable is required")
	}

	supabaseClient, err := supabase.NewClient(databaseURL) 
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}
	defer supabaseClient.Close()

	webhookHandler := &WebhookHandler{
		supabaseClient: supabaseClient,
		webhookSecret: webhookSecret,
	}

	http.HandleFunc("/webhook", webhookHandler.handleWebhook)
	http.HandleFunc("/health", healthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

    log.Printf("Webhook server running on port %s", port)
    log.Printf("Webhook endpoint: http://localhost:%s/webhook", port)
    log.Printf("Health check: http://localhost:%s/health", port)

    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func (h *WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		http.Error(w, "Missing signature header", http.StatusBadRequest)
		return
	}

	if !h.verifySignature(body, signature) {
		log.Printf("Invalid signature: %s", signature)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	githubEvent := r.Header.Get("X-GitHub-Event")
	if githubEvent == "" {
		http.Error(w, "Missing GitHub event header", http.StatusBadRequest)
		return
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Failed to unmarshal payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

    switch githubEvent {
    case "ping":
        log.Println("Ping event received - webhook is working!")

    case "push":
        log.Printf("Push event received")
        if payloadMap, ok := payload.(map[string]any); ok {
            if ref, ok := payloadMap["ref"].(string); ok {
                log.Printf("Ref: %s", ref)
            }
            if repo, ok := payloadMap["repository"].(map[string]any); ok {
                if fullName, ok := repo["full_name"].(string); ok {
                    log.Printf("Repository: %s", fullName)
                }
            }
        }
        // TODO: trigger Cloud Run build for this branch
		// note that there is a "base_ref" field that indiciates main branch
		// so, we can simply check if the incoming ref is main or not; no need to store the main branch in database

	case "create":
		log.Println("Branch or tag creation event received")
		if payloadMap, ok := payload.(map[string]any); ok {
			if ref, ok := payloadMap["ref"].(string); ok {
				log.Printf("Ref created: %s", ref)
			}
			if refType, ok := payloadMap["ref_type"].(string); ok {
				log.Printf("Ref type: %s", refType)
			}
			if repo, ok := payloadMap["repository"].(map[string]any); ok {
				if fullName, ok := repo["full_name"].(string); ok {
					log.Printf("Repository: %s", fullName)
				}
			}
		}
		// TODO: handle branch or tag creation
		// note that GitHub also sends a 'push' event for branch creation, but it's not guaranteed that the 'create' event is first
		// so, the 'push' handler should actually handle the branch creation. 'create' is then just for logging purposes

	case "delete":
		log.Println("Branch or tag deletion event received")
		if payloadMap, ok := payload.(map[string]any); ok {
			if ref, ok := payloadMap["ref"].(string); ok {
				log.Printf("Ref deleted: %s", ref)
			}
			if refType, ok := payloadMap["ref_type"].(string); ok {
				log.Printf("Ref type: %s", refType)
			}
			if repo, ok := payloadMap["repository"].(map[string]any); ok {
				if fullName, ok := repo["full_name"].(string); ok {
					log.Printf("Repository: %s", fullName)
				}
			}
		}
		// TODO: handle branch or tag deletion

    case "installation":
        log.Println("Installation event")
        if err := events.HandleInstallationEvent(h.supabaseClient, payload); err != nil {
			log.Printf("Failed to handle installation event: %v", err)
			return
		}

    case "installation_repositories":
        log.Println("Installation repositories event")
        if err := events.HandleInstallationRepositoriesEvent(h.supabaseClient, payload); err != nil {
			log.Printf("Failed to handle installation repositories event: %v", err)
			return
		}

    default:
        log.Printf("Unhandled event: %s", githubEvent)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *WebhookHandler) verifySignature(payload []byte, signature string) bool {
    mac := hmac.New(sha256.New, []byte(h.webhookSecret))
    mac.Write(payload)
    expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
