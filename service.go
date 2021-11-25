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
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			handleGet(s, w, req)
		case "POST":
			handlePost(s, w, req)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
	})

	return http.ListenAndServe(s.addr, nil)
}

type Document struct {
	ID     string
	Source string
}

type GetResponseBody struct {
	Documents []Document
}

// Get takes the tokens in the query, intersects them and returns a list of
// the matching documents.
func handleGet(s *service, w http.ResponseWriter, req *http.Request) {
	tokens := strings.Split(req.URL.Query().Get("tokens"), ",")
	if len(tokens) == 0 {
		log.Printf("no tokens provided")
		http.Error(w, "", http.StatusBadRequest)
	}

	ids, err := s.querier.Intersect(tokens...)
	if err != nil {
		log.Printf("intersect: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var docs []Document
	for _, id := range ids {
		source, err := s.store.Get(id)
		if err != nil {
			log.Printf("get: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		docs = append(docs, Document{
			ID:     id,
			Source: string(source),
		})
	}

	jsonResp, err := json.Marshal(GetResponseBody{
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
func handlePost(s *service, w http.ResponseWriter, req *http.Request) {
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
