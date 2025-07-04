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

	// handle file uploads, if any
    var attachmentURLs []string
    files := r.MultipartForm.File["files"]
    
    if len(files) > 0 {
        log.Printf("Processing %d file(s)", len(files))
        
        for _, fileHeader := range files {
            file, err := fileHeader.Open()
            if err != nil {
                log.Printf("Error opening file %s: %v", fileHeader.Filename, err)
                continue
            }
            defer file.Close()

            // TODO: Upload to S3
            // but for now, just log the file info
            log.Printf("File: %s, Size: %d bytes, Type: %s", 
                fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
            
            // placeholder URL - replace with actual S3 upload
            attachmentURL := "https://placeholder-s3-url.com/" + fileHeader.Filename
            attachmentURLs = append(attachmentURLs, attachmentURL)
        }
    }

	// postgres expects attachment URLs as a JSON array, so marshal it
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

	// generate embeddings
	// FUTURE - this should be background; consistently takes 800-1000ms
	embedding, err := utils.GetThoughtEmbedding(context.Background(), h.geminiClient, thoughtText)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		http.Error(w, "Failed to generate embedding", http.StatusInternalServerError)
		return
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
	
	// 3. TODO - delete the files from S3 if they exist (database call will return attachment URLs)

	// for now, just log attachment URLs and return success
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

