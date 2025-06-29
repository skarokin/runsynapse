package handlers

import (
	"net/http"
	"log"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/skarokin/runsynapse/go/types"
	"github.com/skarokin/runsynapse/go/utils"
)

func (h *Handler) searchThoughts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request types.SearchThoughtsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := request.UserID
	userID, err := uuid.Parse(string(userIDStr))
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	query := request.Query
	if query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	queryStr := string(query)
	log.Printf("[SEARCH] Searching thoughts for user %s with query: %s", userID, queryStr)

	// get query embedding
	embedding, err := utils.GetQueryEmbedding(r.Context(), h.geminiClient, queryStr)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		http.Error(w, "Failed to generate query embedding", http.StatusInternalServerError)
		return
	}
	// temp print for compile
	log.Printf("[SEARCH] Generated embedding for query: %v", embedding)

	// db call performs hybrid search: does FTS & vector search

	// rerank results, give more weight to FTS results

	// call gemini w/ reranked results
	
	// return

}