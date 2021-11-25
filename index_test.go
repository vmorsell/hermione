package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextDocID(t *testing.T) {
	idx := NewIndex().(*index)

	for i := 0; i < 10; i++ {
		id := idx.nextDocID()
		require.Equal(t, i, id)
	}
}

func TestIndexDocument(t *testing.T) {
	tests := []struct {
		name         string
		dict         map[string][]int
		docIDCounter int
		r            io.Reader
		wantDict     map[string][]int
		err          error
	}{
		{
			name: "ok, first document",
			dict: map[string][]int{},
			r:    strings.NewReader("Hello, world!"),
			wantDict: map[string][]int{
				"hello": {0},
				"world": {0},
			},
		},
		{
			name: "ok, second document",
			dict: map[string][]int{
				"hello": {0},
			},
			docIDCounter: 1,
			r:            strings.NewReader("Hello, world!"),
			wantDict: map[string][]int{
				"hello": {0, 1},
				"world": {1},
			},
		},
	}

	for _, tt := range tests {
		idx := NewIndex().(*index)
		idx.dict = tt.dict
		idx.docIDCounter = tt.docIDCounter

		idx.IndexDocument(tt.r)
		require.EqualValues(t, tt.wantDict, idx.dict)
	}
}

func TestGetPostingsList(t *testing.T) {
	tests := []struct {
		name  string
		dict  map[string][]int
		token string
		res   []int
		err   error
	}{
		{
			name: "not ok, token not found in dict",
			dict: map[string][]int{
				"other": {1},
			},
			token: "x",
			err:   errTokenNotInIndex("x"),
		},
		{
			name: "ok",
			dict: map[string][]int{
				"x": {1},
			},
			token: "x",
			res:   []int{1},
		},
	}

	for _, tt := range tests {
		idx := NewIndex().(*index)
		idx.dict = tt.dict

		res, err := idx.Postings(tt.token)
		require.Equal(t, tt.err, err)
		require.Equal(t, tt.res, res)
	}
}
