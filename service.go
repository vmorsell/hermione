package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Service interface {
	Start() error
}

type service struct {
	addr    string
	idx     Index
	querier Querier
	store   Store
}

func NewService(idx Index, querier Querier, store Store) Service {
	return &service{
		addr:    ":5001",
		idx:     idx,
		querier: querier,
		store:   store,
	}
}

func (s *service) Start() error {
	http.HandleFunc("/search/intersection", s.handleIntersectionSearch)
	http.HandleFunc("/search/phrase", s.handlePhraseSearch)
	http.HandleFunc("/doc", s.handleDoc)

	http.HandleFunc("/debug/postings", s.handleDebugPostings)

	return http.ListenAndServe(s.addr, nil)
}

type Document struct {
	ID     int
	Source string
}

type GetResponseBody struct {
	Hits      int
	Documents []Document
}

// handleIntersectionSearch takes a search query and returns the matching documents
// using intersect.
func (s *service) handleIntersectionSearch(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		log.Printf("unsupported http method: %s", req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query().Get("query")
	if len(query) == 0 {
		log.Printf("no query provided")
		http.Error(w, "", http.StatusBadRequest)
	}

	tokens := strings.Split(query, " ")
	postings, err := s.querier.Intersection(tokens...)
	if err != nil {
		log.Printf("intersection: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var docs []Document
	for _, p := range postings {
		source, err := s.store.Get(p.DocID)
		if err != nil {
			log.Printf("get: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		docs = append(docs, Document{
			ID:     p.DocID,
			Source: string(source),
		})
	}

	jsonResp, err := json.Marshal(GetResponseBody{
		Hits:      len(docs),
		Documents: docs,
	})
	if err != nil {
		log.Printf("marshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

// handlePhraseSearch search for an exact phrase and returns the matching documents.
func (s *service) handlePhraseSearch(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		log.Printf("unsupported http method: %s", req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query().Get("query")
	if len(query) == 0 {
		log.Printf("no query provided")
		http.Error(w, "", http.StatusBadRequest)
	}

	postings, err := s.querier.Phrase(query)
	if err != nil {
		log.Printf("phrase: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var docs []Document
	for _, p := range postings {
		source, err := s.store.Get(p.DocID)
		if err != nil {
			log.Printf("get: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		docs = append(docs, Document{
			ID:     p.DocID,
			Source: string(source),
		})
	}

	jsonResp, err := json.Marshal(GetResponseBody{
		Hits:      len(docs),
		Documents: docs,
	})
	if err != nil {
		log.Printf("marshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

// handlePost takes an document in the body, indexes it and stores it to disk.
func (s *service) handleDoc(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("unsupported http method: %s", req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	buf := bytes.Buffer{}
	r := io.TeeReader(req.Body, &buf)

	id, err := s.idx.IndexDocument(r)
	if err != nil {
		log.Printf("index document: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	r2 := bytes.NewReader(buf.Bytes())
	err = s.store.PutFromStream(r2, id)
	if err != nil {
		log.Printf("put from stream: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

type PostingsBody struct {
	Len       int
	Documents []Posting
}

// handleDebugPostings serves a full postings list.
// The query accepts a token as a query parameter.
func (s *service) handleDebugPostings(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		log.Printf("unsupported http method: %s", req.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	token := req.URL.Query().Get("token")
	if len(token) == 0 {
		log.Printf("no token provided")
		http.Error(w, "", http.StatusBadRequest)
	}

	postings, err := s.idx.Postings(token)
	if err != nil {
		log.Printf("postings: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(PostingsBody{
		Len:       len(postings),
		Documents: postings,
	})
	if err != nil {
		log.Printf("marshal: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}
