package index

import (
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"the quick brown fox", []string{"quick", "brown", "fox"}},
		{"Running runs runner", []string{"run", "run", "runner"}},
		{"", nil},
	}

	for _, tc := range tests {
		got := Analyze(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("Analyze(%q) = %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("Analyze(%q)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}

func TestSearch(t *testing.T) {
	idx := New()
	idx.Add(Document{ID: 1, Title: "Go programming", Body: "Go is a statically typed compiled language"})
	idx.Add(Document{ID: 2, Title: "Python basics", Body: "Python is a dynamically typed interpreted language"})
	idx.Add(Document{ID: 3, Title: "Go web servers", Body: "Go is great for building fast web servers"})
	idx.Finalize()

	tests := []struct {
		query   string
		wantIDs []uint32
	}{
		{"go language", []uint32{1}},
		{"web servers", []uint32{3}},
		{"python", []uint32{2}},
		{"java", nil},
		{"", nil},
	}

	for _, tc := range tests {
		results := idx.Search(tc.query)
		if len(results) != len(tc.wantIDs) {
			t.Errorf("Search(%q) returned %d results, want %d", tc.query, len(results), len(tc.wantIDs))
			continue
		}
		for i, doc := range results {
			if doc.ID != tc.wantIDs[i] {
				t.Errorf("Search(%q)[%d].ID = %d, want %d", tc.query, i, doc.ID, tc.wantIDs[i])
			}
		}
	}
}
