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
		HasMoreAbove: false,
		HasMoreBelow: false,
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
	userID := request.UserID
	limit := request.Limit
	cursor := request.Cursor
	order := request.Order

	log.Println("[LOAD] Loading thoughts for user:", userID, "with limit:", limit, "and cursor:", cursor, "order:", order)

	// infinite scrolling logic:
	//    - if no cursor, load the latest [limit] thoughts
	//    - if cursor is provided...
	//        - if direction is "before", load [limit] thoughts before the cursor timestamp
	//        - if direction is "after", load [limit] thoughts after the cursor timestamp
	// 4. return the newly loaded thoughts && whether or not there are more thoughts above or below the current cursor 

	return 
}