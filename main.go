package main

import (
	"log"
)

func main() {
	idx := NewIndex()
	querier := NewQuerier(idx)
	store, err := NewStore("./store")
	if err != nil {
		log.Fatalf("new store: %v", err)
	}

	s := NewService(idx, querier, store)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
