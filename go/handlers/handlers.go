package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/genai"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Handler struct {
	supabaseClient *pgxpool.Pool
	geminiClient   *genai.Client
	s3Client 	   *s3.Client
	s3Bucket       string
	mux            *http.ServeMux
}

// upon registering a new handler, setup routes
func NewHandler(supabase *pgxpool.Pool, gemini *genai.Client, s3 *s3.Client, s3Bucket string) *Handler {
	h := &Handler{
		supabaseClient: supabase,
		geminiClient:   gemini,
		s3Client: 	 	s3,
		s3Bucket:       s3Bucket,
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