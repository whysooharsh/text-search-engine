package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/whysooharsh/text-search-engine/index"
)

type Server struct {
	idx *index.Index
	mux *http.ServeMux
}

type searchResult struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}

type searchResponse struct {
	Query     string         `json:"query"`
	Total     int            `json:"total"`
	TimeTaken string         `json:"time_taken"`
	Results   []searchResult `json:"results"`
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
		http.Error(w, `{"error":"missing q param"}`, http.StatusBadRequest)
		return
	}

	start := time.Now()
	docs := s.idx.Search(query)
	elapsed := time.Since(start)

	log.Printf("SEARCH | query: %-20q | results: %4d | time: %s", query, len(docs), elapsed)

	results := make([]searchResult, len(docs))
	for i, doc := range docs {
		results[i] = searchResult{ID: doc.ID, Title: doc.Title}
	}

	resp := searchResponse{
		Query:     query,
		Total:     len(docs),
		TimeTaken: elapsed.String(),
		Results:   results,
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
}
