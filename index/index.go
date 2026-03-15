package index

import (
	"sort"
	"sync"
)

// Document represents a single searchable document.
type Document struct {
	ID    uint32
	Title string
	Body  string
}

// Index is a thread-safe inverted index that maps terms to document IDs.
type Index struct {
	mu       sync.RWMutex
	postings map[string][]uint32
	docs     map[uint32]Document
}

// New creates and returns an empty Index.
func New() *Index {
	return &Index{
		postings: make(map[string][]uint32),
		docs:     make(map[uint32]Document),
	}
}

// Add indexes a single document. It is safe to call Add concurrently.
func (idx *Index) Add(doc Document) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.docs[doc.ID] = doc

	seen := make(map[string]struct{})
	for _, term := range Analyze(doc.Title + " " + doc.Body) {
		if _, dup := seen[term]; dup {
			continue
		}
		seen[term] = struct{}{}
		idx.postings[term] = append(idx.postings[term], doc.ID)
	}
}

// Finalize sorts every posting list so that subsequent searches can use
// the two-pointer intersection algorithm. Call once after all Add calls.
func (idx *Index) Finalize() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for term := range idx.postings {
		sort.Slice(idx.postings[term], func(i, j int) bool {
			return idx.postings[term][i] < idx.postings[term][j]
		})
	}
}

// Postings returns the sorted list of document IDs that contain term.
func (idx *Index) Postings(term string) []uint32 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.postings[term]
}

// GetDoc returns the Document with the given id and whether it was found.
func (idx *Index) GetDoc(id uint32) (Document, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	doc, ok := idx.docs[id]
	return doc, ok
}
