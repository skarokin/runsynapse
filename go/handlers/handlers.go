package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/genai"
)

type Handler struct {
	supabaseClient *pgxpool.Pool
	geminiClient   *genai.Client
}

func NewHandler(supabase *pgxpool.Pool, gemini *genai.Client) *Handler {
	return &Handler{
		supabaseClient: supabase,
		geminiClient:   gemini,
	}
}

func (h *Handler) SetupRoutes() {
	http.HandleFunc("/loadFunction", h.loadFunction)
	http.HandleFunc("/loadThoughts", h.loadThoughts)
	http.HandleFunc("/pinThought", h.pinThought)
	http.HandleFunc("/unpinThought", h.unpinThought)
	http.HandleFunc("/searchThoughts", h.searchThoughts)
	http.HandleFunc("/gotoPin", h.gotoPin)
	http.HandleFunc("/deleteThought", h.deleteThought)
	http.HandleFunc("/newThought", h.newThought)
	http.HandleFunc("/health", h.healthCheck)
}
