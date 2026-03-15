package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "doc1.txt"), []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "doc2.txt"), []byte("foo bar"), 0o644); err != nil {
		t.Fatal(err)
	}

	docs, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("expected 2 docs, got %d", len(docs))
	}
}

func TestLoad_IgnoresNonTxt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "doc.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "readme.md"), []byte("markdown"), 0o644); err != nil {
		t.Fatal(err)
	}

	docs, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Errorf("expected 1 doc, got %d", len(docs))
	}
}

func TestLoad_IgnoresSubdirectories(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "doc.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	docs, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Errorf("expected 1 doc, got %d", len(docs))
	}
}

func TestLoad_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for empty directory, got nil")
	}
}

func TestLoad_MissingDirectory(t *testing.T) {
	_, err := Load("/tmp/nonexistent-dir-that-does-not-exist-xyz")
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

func TestLoad_TitleStripsExtension(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "my-document.txt"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	docs, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
	if docs[0].Title != "my-document" {
		t.Errorf("expected title %q, got %q", "my-document", docs[0].Title)
	}
}

func TestLoad_DocBody(t *testing.T) {
	dir := t.TempDir()
	body := "the quick brown fox"
	if err := os.WriteFile(filepath.Join(dir, "fox.txt"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	docs, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docs[0].Body != body {
		t.Errorf("expected body %q, got %q", body, docs[0].Body)
	}
}

func TestLoad_IDsStartAtOne(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	docs, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docs[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", docs[0].ID)
	}
}
