package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/skarokin/runsynapse/go/types"
)

// exact same as loadThoughts but returns the pinned thoughts too
func (h* Handler) loadFunction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request types.RequestsOnlyRequiringUserID
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// extract user ID (it's of type UserID, already validated)
	userIDStr := request.UserID
	userID, err := uuid.Parse(string(userIDStr))
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	log.Println("[LOAD] Load function called for user:", userIDStr)

	// get db result
	var res string
    err = h.supabaseClient.QueryRow(context.Background(), `
        SELECT load_thoughts_and_pins($1)
    `, userID).Scan(&res)

    if err != nil {
        log.Printf("Error loading thoughts: %v", err)
        http.Error(w, "Failed to load thoughts", http.StatusInternalServerError)
        return
    }

	// parse result into a struct
    var dbResult struct {
        Thoughts       []json.RawMessage `json:"thoughts"`
        PinnedThoughts []json.RawMessage `json:"pinned_thoughts"`
		HasMoreAbove   bool              `json:"has_more_above"`
    }

	// the result is a JSON string, so we need to unmarshal it
    err = json.Unmarshal([]byte(res), &dbResult)
    if err != nil {
        log.Printf("Error parsing database result: %v", err)
        http.Error(w, "Failed to parse result", http.StatusInternalServerError)
        return
    }

    // convert raw messages to Thought structs
    var thoughts []types.Thought
    for _, rawThought := range dbResult.Thoughts {
        var thought types.Thought
        if err := json.Unmarshal(rawThought, &thought); err != nil {
            log.Printf("Error unmarshaling thought: %v", err)
            continue
        }
        thoughts = append(thoughts, thought)
    }

    var pinnedThoughts []types.Thought
    for _, rawThought := range dbResult.PinnedThoughts {
        var thought types.Thought
        if err := json.Unmarshal(rawThought, &thought); err != nil {
            log.Printf("Error unmarshaling pinned thought: %v", err)
            continue
        }
        pinnedThoughts = append(pinnedThoughts, thought)
    }

	// build response
	response := types.LoadFunctionResponse{
		Thoughts:       thoughts,
		PinnedThoughts: pinnedThoughts,
		HasMoreAbove: dbResult.HasMoreAbove,
	}

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *Handler) loadThoughts(w http.ResponseWriter, r *http.Request) {
	var request types.LoadThoughtsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// extract userID and params from request
	userIDStr := request.UserID
	cursorStr := request.Cursor

	log.Println("[LOAD] Loading thoughts for user:", userIDStr, "with cursor:", cursorStr)

	userID, err := uuid.Parse(string(userIDStr))
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	cursor, err := uuid.Parse(string(cursorStr))
	if err != nil && cursorStr != "" {
		log.Printf("Invalid cursor: %v", err)
		http.Error(w, "Invalid cursor", http.StatusBadRequest)
		return
	}

	// get db result
	var res string
	err = h.supabaseClient.QueryRow(context.Background(), `
		SELECT load_more($1, $2)
	`, userID, cursor).Scan(&res)

	if err != nil {
		log.Printf("Error loading thoughts: %v", err)
		http.Error(w, "Failed to load thoughts", http.StatusInternalServerError)
		return
	}

	// parse result into a struct
	var dbResult struct {
		Thoughts       []json.RawMessage  `json:"thoughts"`
		HasMoreAbove   bool               `json:"has_more_above"`
	}
	err = json.Unmarshal([]byte(res), &dbResult)
	if err != nil {
		log.Printf("Error parsing database result: %v", err)
		http.Error(w, "Failed to parse result", http.StatusInternalServerError)
		return
	}

	// convert raw messages to Thought structs
	var thoughts []types.Thought
	for _, rawThought := range dbResult.Thoughts {
		var thought types.Thought
		if err := json.Unmarshal(rawThought, &thought); err != nil {
			log.Printf("Error unmarshaling thought: %v", err)
			continue
		}
		thoughts = append(thoughts, thought)
	}

	// build response
	response := types.LoadThoughtsResponse{
		Thoughts:       thoughts,
		HasMoreAbove: dbResult.HasMoreAbove,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}