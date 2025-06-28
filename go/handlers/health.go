package handlers

import (
	"net/http"
	"log"
)

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("[HEALTH] Health check")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
