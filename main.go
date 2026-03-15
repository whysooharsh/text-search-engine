package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown failed: %v", err)
		}
	}()

	log.Println("listening on:", *addr)
	if err := server.ListenAndServe(*addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
	log.Println("server stopped")
}
