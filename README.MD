# Text Search Engine

A lightweight full-text search engine in Go. Indexes `.txt` documents and serves search queries over HTTP.

## How it works

1. Loads `.txt` files from a directory
2. Tokenizes, removes stop words, and stems each document
3. Builds an inverted index mapping terms to document IDs
4. Serves a search API using two-pointer posting list intersection

## Project structure

```
search-engine/
├── main.go           entry point, CLI flags
├── loader/
│   └── loader.go     reads .txt files from a directory
├── index/
│   ├── analyze.go    tokenize, stopword filter, stem
│   ├── index.go      inverted index core
│   └── search.go     boolean AND search via two-pointer intersection
├── api/
│   └── api.go        HTTP search endpoint
└── doc/              sample documents for local testing
```

## Getting started

**Prerequisites:** Go 1.21+

```bash
go mod tidy
go run . -docs ./doc
```

```bash
curl "http://localhost:8000/search?q=machine%20learning"
```

## API

### `GET /search?q=<query>`

| Parameter | Type   | Description        |
|-----------|--------|--------------------|
| `q`       | string | search query terms |

**Response**
```json
{
  "query": "machine learning",
  "total": 1,
  "results": [{ "id": 2, "title": "doc2" }]
}
```

## CLI flags

| Flag    | Default | Description                      |
|---------|---------|----------------------------------|
| `-docs` | `./doc` | path to folder with .txt files   |
| `-addr` | `:8000` | HTTP address to listen on        |

## Tests

```bash
go test ./...
```

## Built with

- [snowball](https://github.com/kljensen/snowball) - Go implementation of the Snowball stemmer