package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

func main() {
	idx := NewIndex()
	querier := NewQuerier(idx)
	store, err := NewStore("./store")
	if err != nil {
		log.Fatalf("new store: %v", err)
	}

	go func() {
		s := NewService(idx, querier, store)
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Add four documents to the index.
	docs := []string{
		"new home sales top forecasts",
		"home sales rise in july",
		"increase in home sales in july",
		"july new home sales rise",
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	for _, d := range docs {
		r := bytes.NewReader([]byte(d))
		req, err := http.NewRequest("POST", "http://localhost:5001", r)
		if err != nil {
			log.Printf("new request: %v", err)
		}
		req.Header.Set("Content-Type", "application/octet-stream")

		res, err := client.Do(req)
		if err != nil {
			log.Printf("do req: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("failed with status: %d", res.StatusCode)
		}
	}

	// Block main thread.
	select {}
}
