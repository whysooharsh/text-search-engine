package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/whysooharsh/text-search-engine/index"
)

// Server holds the index and routes all HTTP requests.
type Server struct {
	idx    *index.Index
	mux    *http.ServeMux
	server *http.Server
}

type searchResponse struct {
	Query   string          `json:"query"`
	Total   int             `json:"total"`
	Results []searchResults `json:"results"`
}

type searchResults struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}

// New creates a Server wired up to idx and registers all routes.
func New(idx *index.Index) *Server {
	s := &Server{idx: idx, mux: http.NewServeMux()}
	s.server = &http.Server{Handler: s.mux}
	s.mux.HandleFunc("/search", s.handleSearch)
	s.mux.HandleFunc("/health", s.handleHealth)
	return s
}

// ListenAndServe starts the HTTP server on addr.
func (s *Server) ListenAndServe(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

// Shutdown gracefully drains connections using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing query param"})
		return
	}

	docs := s.idx.Search(query)

	resp := searchResponse{
		Query:   query,
		Total:   len(docs),
		Results: make([]searchResults, len(docs)),
	}
	for i, doc := range docs {
		resp.Results[i] = searchResults{ID: doc.ID, Title: doc.Title}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
