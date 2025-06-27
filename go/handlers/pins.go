package handlers

import (
	"net/http"
)

func (h* Handler) pinThought(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. parse request json
	// 2. update the pins table for this user in the database
	// 3. return the updated list of pinned thoughts for immediate use in the UI
}

func (h *Handler) unpinThought(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. parse request json
	// 2. remove the pin from the pins table for this user in the database
	// 3. return the updated list of pinned thoughts for immediate use in the UI
}

func (h *Handler) gotoPin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. parse request json
	// 2. extract the pinned thought ID
	// 3. perform a 'load thoughts' operation using the pinned thought ID's timestamp as the cursor
	// 4. return the thoughts loaded around the pinned thought
}