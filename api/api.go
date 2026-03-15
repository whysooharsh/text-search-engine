package api

import (
	"encoding/json"
	"net/http"

	"github.com/whysooharsh/text-search-engine/index"
)

type Server struct {
	idx *index.Index
	mux *http.ServeMux
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

func New(idx *index.Index) *Server {

	s := &Server{idx: idx, mux: http.NewServeMux()}
	s.mux.HandleFunc("/search", s.handleSearch)

	return s
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, `{"error":"missing query params"}`, http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
