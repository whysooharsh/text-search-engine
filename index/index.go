package index

import (
	"sort"
	"sync"
)

type Document struct {
	ID    uint32
	Title string
	Body  string
}

type Index struct {
	mu       sync.RWMutex
	postings map[string][]uint32
	docs     map[uint32]Document
}

func New() *Index {
	return &Index{
		postings: make(map[string][]uint32),
		docs:     make(map[uint32]Document),
	}
}

func (idx *Index) Add(doc Document) {
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

func (idx *Index) Finalize() {
	for term := range idx.postings {
		sort.Slice(idx.postings[term], func(i, j int) bool {
			return idx.postings[term][i] < idx.postings[term][j]
		})
	}
}

func (idx *Index) Postings(term string) []uint32 {

	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.postings[term]
}
func (idx *Index) GetDoc(id uint32) (Document, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	doc, ok := idx.docs[id]
	return doc, ok
}
