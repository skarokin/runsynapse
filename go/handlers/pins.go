package handlers

import (
	"net/http"
	"log"
	"encoding/json"

	"github.com/skarokin/runsynapse/go/types"
)

func (h *Handler) pinThought(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request types.TogglePinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("[PINS] Pinning thought for user:", request.UserID)

	// 1. parse request json
	// 2. update the pins table for this user in the database
	// 3. return the updated list of pinned thoughts for immediate use in the UI
}

func (h *Handler) unpinThought(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request types.TogglePinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("[PINS] Unpinning thought for user:", request.UserID)

	// 1. parse request json
	// 2. remove the pin from the pins table for this user in the database
	// 3. return the updated list of pinned thoughts for immediate use in the UI
}