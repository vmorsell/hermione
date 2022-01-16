package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Ingest accepts an arbitrary number of arguments with relative paths to corpuses.
// All provided directories are crawled and all documents are indexed.
func main() {
	// First arg is always the program name.
	if len(os.Args) == 1 {
		log.Fatal("no corpuses provided")
	}

	corpusDirs := os.Args[1:]
	n := 0
	for i, src := range corpusDirs {
		log.Printf("Processing %s (%d/%d)...", src, i+1, len(corpusDirs))
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			log.Printf("%s", path)
			if err := ingestFile(path); err != nil {
				return fmt.Errorf("ingest file: %w", err)
			}

			n++

			return nil
		})
		if err != nil {
			log.Fatalf("walk: %v", err)
		}
		log.Printf("Done processing %s.", src)
	}
	log.Printf("Ingested %d files", n)
}

func ingestFile(src string) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest("POST", "http://localhost:5001/doc", file)
	if err != nil {
		log.Printf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("do req: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("failed with status: %d", res.StatusCode)
	}
	return nil
}
