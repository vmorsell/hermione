package main

import (
	"fmt"
	"io"
	"time"

	"github.com/segmentio/ksuid"
)

type Index interface {
	IndexDocument(r io.Reader) (string, error)
	Postings(token string) ([]string, error)
}

type index struct {
	dict map[string][]string

	idFn func() (string, error)
}

func NewIndex() Index {
	return &index{
		dict: make(map[string][]string),

		// use default functions
		idFn: id,
	}
}

// IndexDocument tokenizes the document from the reader and adds the tokens to
// the index. It returns the ID of the new document.
func (idx *index) IndexDocument(r io.Reader) (string, error) {
	id, err := idx.idFn()
	if err != nil {
		return "", fmt.Errorf("id: %w", err)
	}

	tokenizer := NewTokenizer(r)
	for tokenizer.HasMoreTokens() {
		t, err := tokenizer.NextToken()
		if err != nil {
			return "", fmt.Errorf("next token: %w", err)
		}

		if _, ok := idx.dict[t]; !ok {
			idx.dict[t] = []string{}
		}
		idx.dict[t] = append(idx.dict[t], id)
	}
	return id, nil
}

var errTokenNotInIndex = func(token string) error { return fmt.Errorf("token '%s' not found in index", token) }

// Postings returns the full postings list for the given token.
func (idx *index) Postings(token string) ([]string, error) {
	if _, ok := idx.dict[token]; !ok {
		return nil, errTokenNotInIndex(token)
	}
	return idx.dict[token], nil
}

func id() (string, error) {
	id, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		return "", fmt.Errorf("new random with time: %w", err)
	}
	return id.String(), nil
}
