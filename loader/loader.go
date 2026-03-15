package loader

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/whysooharsh/text-search-engine/index"
)

func Load(dir string) ([]index.Document, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("can't open dir %q : %w", dir, err)
	}

	var docs []index.Document
	var id uint32

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		body, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}

		id++
		docs = append(docs, index.Document{
			ID:    id,
			Title: strings.TrimSuffix(entry.Name(), ".txt"),
			Body:  string(body),
		})
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("no .txt files found in %q", dir)
	}

	return docs, nil

}
