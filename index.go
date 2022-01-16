package main

import (
	"fmt"
	"io"
	"time"

	"github.com/segmentio/ksuid"
)

type Index interface {
	IndexDocument(r io.Reader) (string, error)
	Postings(token string) ([]Posting, error)
}

type index struct {
	dict map[string][]Posting

	idFn func() (string, error)
}

func NewIndex() Index {
	return &index{
		dict: make(map[string][]Posting),

		// use default functions
		idFn: id,
	}
}

// hasDocID searches for a document ID in the postings list, and returns the index if it's found.
func hasDocID(postingsList []Posting, id string) (bool, int) {
	if len(postingsList) == 0 {
		return false, 0
	}

	for i, p := range postingsList {
		if p.DocID > id {
			return false, 0
		}
		if p.DocID == id {
			return true, i
		}
	}
	return false, 0
}

type Posting struct {
	DocID string
	Freq  int
}

// IndexDocument tokenizes the document from the reader and adds the tokens to
// the index. It returns the ID of the new document.
func (idx *index) IndexDocument(r io.Reader) (string, error) {
	id, err := idx.idFn()
	if err != nil {
		return "", fmt.Errorf("id: %w", err)
	}

	tokenizer := NewTokenizer(r)
next:
	for tokenizer.HasMoreTokens() {
		t, err := tokenizer.NextToken()
		if err != nil {
			return "", fmt.Errorf("next token: %w", err)
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

func id() (string, error) {
	id, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		return "", fmt.Errorf("new random with time: %w", err)
	}
	return id.String(), nil
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
