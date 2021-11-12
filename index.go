package main

import (
	"fmt"
	"io"
)

type Index interface {
	IndexDocument(r io.Reader) error
}

type index struct {
	dict         map[string][]int
	docIDCounter int
}

func NewIndex() Index {
	return &index{}
}

// nextDocID returns the next available document ID.
func (idx *index) nextDocID() int {
	idx.docIDCounter++
	return idx.docIDCounter - 1
}

// IndexDocument tokenizes the document from the reader and adds the tokens to
// the index.
func (idx *index) IndexDocument(r io.Reader) error {
	docID := idx.nextDocID()

	tokenizer := NewTokenizer(r)
	for tokenizer.HasMoreTokens() {
		t, err := tokenizer.NextToken()
		if err != nil {
			return fmt.Errorf("next token: %w", err)
		}

		if _, ok := idx.dict[t]; !ok {
			idx.dict[t] = []int{}
		}
		idx.dict[t] = append(idx.dict[t], docID)
	}
	return nil
}

var errTokenNotFound = func(token string) error { return fmt.Errorf("token '%s' not found in index", token) }

func (idx *index) GetPostingsList(token string) ([]int, error) {
	if _, ok := idx.dict[token]; !ok {
		return nil, errTokenNotFound(token)
	}
	return idx.dict[token], nil
}
