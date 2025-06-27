package handlers

import (
	"net/http"
)

func (h *Handler) searchThoughts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. parse request json
	// 2. extract user ID and search query
	// 3. create embedding on search query
	// 4. do basic FTS search AND vector search on user's thoughts (hybrid search)
	//    - limit to 100 results per query
	// 5. merge & rank - any thought that shows up in FTS should receive a significant rank boost
	//    - rank by relevance to query, recency, pin status, and other factors.
	//    - limit to 10 merged results
	// 6. call Gemini with ranked thoughts 
	// 7. return ranked thoughts and Gemini response
}