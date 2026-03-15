package main

import (
	"flag"
	"log"
	"time"

	"github.com/whysooharsh/text-search-engine/api"
	"github.com/whysooharsh/text-search-engine/index"
	"github.com/whysooharsh/text-search-engine/loader"
)

func main() {

	docsDir := flag.String("docs", "./doc", "folder containing text files")
	addr := flag.String("addr", ":8000", "HTTP listen address")
	flag.Parse()

	idx := index.New()

	log.Println("Full text search is in progress")

	start := time.Now()

	docs, err := loader.Load(*docsDir)

	if err != nil {
		log.Fatalf("load: %v", err)
	}

	for _, doc := range docs {
		idx.Add(doc)
	}

	idx.Finalize()

	log.Printf("indexed %d documents in %s\n", len(docs), time.Since(start))
	server := api.New(idx)

	log.Println("listening on : ", *addr)
	log.Fatal(server.ListenAndServe(*addr))

}
