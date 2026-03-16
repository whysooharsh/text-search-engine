package loader

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/whysooharsh/text-search-engine/index"
)

type xmlPage struct {
	Title string `xml:"title"`
	ID    uint32 `xml:"id"`
	Text  string `xml:"revision>text"`
}

func Load(path string, fn func(index.Document), maxDocs int) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cant open file %q: %w", path, err)
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	count := 0

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("xml parse error: %w", err)
		}

		se, ok := tok.(xml.StartElement)
		if !ok || se.Name.Local != "page" {
			continue
		}

		var page xmlPage
		if err := dec.DecodeElement(&page, &se); err != nil {
			continue
		}

		if page.Text == "" || page.Title == "" {
			continue
		}

		fn(index.Document{
			ID:    page.ID,
			Title: page.Title,
			Body:  page.Text,
		})

		count++
		if count%10000 == 0 {
			log.Printf("processed %d documents...", count)
		}
		if maxDocs > 0 && count >= maxDocs {
			log.Printf("reached limit of %d documents, stopping", maxDocs)
			break
		}
	}

	return nil
}
