package main

import (
	"fmt"
	"io"
)

type Index interface {
	IndexDocument(r io.Reader) (int, error)
	Postings(token string) ([]Posting, error)
}

type index struct {
	dict   map[string][]Posting
	nextID int

	idFn func() int
}

func NewIndex() Index {
	return &index{
		dict:   make(map[string][]Posting),
		nextID: 0,
	}
}

// hasDocID searches for a document ID in the postings list, and returns the index if it's found.
func hasDocID(postingsList []Posting, id int) (bool, int) {
	for i, p := range postingsList {
		if p.DocID < id {
			continue
		}
		if p.DocID == id {
			return true, i
		}
		return false, 0
	}
	return false, 0
}

type Posting struct {
	DocID int
	Freq  int
}

// IndexDocument tokenizes the document from the reader and adds the tokens to
// the index. It returns the ID of the new document.
func (idx *index) IndexDocument(r io.Reader) (int, error) {
	id := idx.id()
	tokenizer := NewTokenizer(r)

next:
	for tokenizer.HasMoreTokens() {
		t, err := tokenizer.NextToken()
		if err != nil {
			return 0, fmt.Errorf("next token: %w", err)
		}
		if t == "" {
			break
		}

		if _, ok := idx.dict[t]; !ok {
			idx.dict[t] = []Posting{}
		}

		found, i := hasDocID(idx.dict[t], id)
		if found {
			idx.dict[t][i].Freq++
			continue next
		}

		idx.dict[t] = append(idx.dict[t], Posting{
			DocID: id,
			Freq:  1,
		})
	}
	return id, nil
}

var errTokenNotInIndex = func(token string) error { return fmt.Errorf("token '%s' not found in index", token) }

// Postings returns the full postings list for the given token.
func (idx *index) Postings(token string) ([]Posting, error) {
	if _, ok := idx.dict[token]; !ok {
		return nil, errTokenNotInIndex(token)
	}
	return idx.dict[token], nil
}

func (idx *index) id() int {
	if idx.idFn != nil {
		return idx.idFn()
	}

	id := idx.nextID
	idx.nextID++
	return id
}

func max(n ...int) int {
	if len(n) == 0 {
		return 0
	}

	out := n[0]
	for _, nn := range n[1:] {
		if nn > out {
			out = nn
		}
	}
	return out
}

func min(n ...int) int {
	if len(n) == 0 {
		return 0
	}

	out := n[0]
	for _, nn := range n[1:] {
		if nn < out {
			out = nn
		}
	}
	return out
}
