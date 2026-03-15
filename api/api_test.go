package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whysooharsh/text-search-engine/index"
)

func newTestServer() *Server {
	idx := index.New()
	idx.Add(index.Document{ID: 1, Title: "Go programming", Body: "Go is a statically typed compiled language"})
	idx.Add(index.Document{ID: 2, Title: "Python basics", Body: "Python is a dynamically typed interpreted language"})
	idx.Finalize()
	return New(idx)
}

func TestHandleSearch_MissingQuery(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()

	s.handleSearch(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["error"] == "" {
		t.Error("expected non-empty error field in response")
	}
}

func TestHandleSearch_ValidQuery(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/search?q=go+language", nil)
	rec := httptest.NewRecorder()

	s.handleSearch(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	var resp searchResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if resp.Query != "go language" {
		t.Errorf("expected query %q, got %q", "go language", resp.Query)
	}
	if resp.Total != len(resp.Results) {
		t.Errorf("total %d does not match results length %d", resp.Total, len(resp.Results))
	}
}

func TestHandleSearch_NoResults(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/search?q=rust", nil)
	rec := httptest.NewRecorder()

	s.handleSearch(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var resp searchResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if resp.Total != 0 {
		t.Errorf("expected 0 results, got %d", resp.Total)
	}
}

func TestHandleHealth(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	s.handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %q", body["status"])
	}
}
