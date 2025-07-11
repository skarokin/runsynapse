package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	
	"github.com/skarokin/runsynapse/go/types"
	"github.com/skarokin/runsynapse/go/utils"
)

type FileResult struct {
	URLs []string
	Err error
}

type EmbeddingResult struct {
	Embedding string
	Err error
}

func (h *Handler) newThought(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse mulipart form
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

    userIDstr := r.FormValue("user_id")
    thoughtText := r.FormValue("thought")

    if userIDstr == "" {
		log.Println("user_id is required")
        http.Error(w, "user_id is required", http.StatusBadRequest)
        return
    }

	log.Println("[NEW_THOUGHT] New thought request received for user:", userIDstr)

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

    if thoughtText == "" {
		log.Println("thought is required")
        http.Error(w, "thought is required", http.StatusBadRequest)
        return
    }

	// make channels for files and embedding async calls
	fileChan := make(chan FileResult, 1)
	embeddingChan := make(chan EmbeddingResult, 1)
    
	// process attachments
    go func() {
        defer close(fileChan)
        
        urls, err := utils.UploadFilesToS3(
            context.Background(),
            h.s3Client,
            h.s3Bucket,
            r.MultipartForm.File["files"],
        )
        fileChan <- FileResult{URLs: urls, Err: err}
    }()

	// generate embeddings
    go func() {
        defer close(embeddingChan)
        
        embedding, err := utils.GetThoughtEmbedding(context.Background(), h.geminiClient, thoughtText)
        embeddingChan <- EmbeddingResult{Embedding: embedding, Err: err}
    }()

	// wait for both goroutines to finish
	fileResult := <-fileChan
	embeddingResult := <-embeddingChan

	if fileResult.Err != nil {
		log.Printf("Error uploading files: %v", fileResult.Err)
		http.Error(w, "Failed to upload files", http.StatusInternalServerError)
		return
	}

	if embeddingResult.Err != nil {
		log.Printf("Error generating embedding: %v", embeddingResult.Err)
		http.Error(w, "Failed to generate embedding", http.StatusInternalServerError)
		return
	}

	attachmentURLs := fileResult.URLs
	embedding := embeddingResult.Embedding

	// postgres expects attachment URLs as a JSON array so marshal it
	var attachmentURLsBytes []byte
    if len(attachmentURLs) == 0 {
        attachmentURLsBytes = []byte("[]")
    } else {
        attachmentURLsBytes, err = json.Marshal(attachmentURLs)
        if err != nil {
            log.Printf("Error marshalling attachment URLs: %v", err)
            http.Error(w, "Failed to process attachments", http.StatusInternalServerError)
            return
        }
    }

	// get db result
	var res string

    err = h.supabaseClient.QueryRow(context.Background(), `
        SELECT * FROM new_thought($1, $2, $3, $4)
    `, userID, thoughtText, embedding, string(attachmentURLsBytes)).Scan(&res)
	if err != nil {
		log.Printf("Error inserting new thought: %v", err)
		http.Error(w, "Failed to insert new thought", http.StatusInternalServerError)
		return
	}

	// parse result into a struct
	var dbResult struct {
		ID string `json:"id"`
		CreatedAt string `json:"created_at"`
	}

	// the result is a JSON string, so we need to unmarshal it
	err = json.Unmarshal([]byte(res), &dbResult)
    if err != nil {
        log.Printf("Error parsing database result: %v", err)
        http.Error(w, "Failed to parse result", http.StatusInternalServerError)
        return
    }

	// validate the thought ID
    thoughtID, err := uuid.Parse(dbResult.ID); if err != nil {
		log.Printf("Invalid thought ID from database: %v", err)
		http.Error(w, "Invalid thought ID", http.StatusInternalServerError)
		return
	}

	// build the resposne
	newThought := types.Thought{
		ID:        thoughtID,
		Thought:   thoughtText,
		Pinned:    false, // default to not pinned
		Created: dbResult.CreatedAt,
		Attachments: attachmentURLs,
	}

    response := types.NewThoughtResponse{
        Thought: newThought,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *Handler) deleteThought(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed for deleteThought")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request types.DeleteThoughtRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userIDStr := request.UserID
	thoughtIDStr := request.ThoughtID

	if thoughtIDStr == "" {
		log.Println("thought_id is required")
		http.Error(w, "thought_id is required", http.StatusBadRequest)
		return
	}

	if userIDStr == "" {
		log.Println("user_id is required")
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	log.Println("[DELETE_THOUGHT] Delete thought request received for user:", userIDStr)

	userID, err := uuid.Parse(string(userIDStr))
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	thoughtID, err := uuid.Parse(string(thoughtIDStr))
	if err != nil {
		log.Printf("Invalid thought_id: %v", err)
		http.Error(w, "Invalid thought_id", http.StatusBadRequest)
		return
	}

	// 2. delete the thought and associated data from the database
	var res string
	err = h.supabaseClient.QueryRow(context.Background(), `
		SELECT * FROM delete_thought($1, $2)
	`, userID, thoughtID).Scan(&res)
	if err != nil {
		log.Printf("Error deleting thought: %v", err)
		http.Error(w, "Failed to delete thought", http.StatusInternalServerError)
		return
	}

	// parse result into a struct
	var dbResult struct {
		Deleted bool `json:"deleted"`
		AttachmentURLs []string `json:"attachment_urls"`
		ThoughtID string `json:"thought_id"`
	}

	err = json.Unmarshal([]byte(res), &dbResult)
	if err != nil {
		log.Printf("Error parsing database result: %v", err)
		http.Error(w, "Failed to parse result", http.StatusInternalServerError)
		return
	}
	
	// 3. delete the files from S3 if they exist (database call will return attachment URLs)
	for _, url := range dbResult.AttachmentURLs {
		if url == "" {
			continue // skip empty URLs
		}
		err = utils.DeleteFromS3(
			context.Background(),
			h.s3Client,
			h.s3Bucket,
			url,
		)
		if err != nil {
			log.Printf("Error deleting file from S3: %v", err)
		}
	}

	response := types.DeleteThoughtResponse{
		Success: dbResult.Deleted,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}