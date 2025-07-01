package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/genai"
)

type Handler struct {
	supabaseClient *pgxpool.Pool
	geminiClient   *genai.Client
	mux            *http.ServeMux
}

func NewHandler(supabase *pgxpool.Pool, gemini *genai.Client) *Handler {
	h := &Handler{
		supabaseClient: supabase,
		geminiClient:   gemini,
		mux:            http.NewServeMux(),
	}
	h.setupRoutes()
	return h
}

// registers all routes (private method)
func (h *Handler) setupRoutes() {
	h.mux.HandleFunc("/loadFunction", h.loadFunction)
	h.mux.HandleFunc("/loadThoughts", h.loadThoughts)
	h.mux.HandleFunc("/pinThought", h.pinThought)
	h.mux.HandleFunc("/unpinThought", h.unpinThought)
	h.mux.HandleFunc("/searchThoughts", h.searchThoughts)
	h.mux.HandleFunc("/deleteThought", h.deleteThought)
	h.mux.HandleFunc("/newThought", h.newThought)
	h.mux.HandleFunc("/health", h.healthCheck)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}