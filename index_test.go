package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextDocID(t *testing.T) {
	idx := index{}

	for i := 0; i < 10; i++ {
		id := idx.nextDocID()
		require.Equal(t, i, id)
	}
}

func TestIndexDocument(t *testing.T) {
	tests := []struct {
		name     string
		idx      *index
		r        io.Reader
		wantDict map[string][]int
		err      error
	}{
		{
			name: "ok, first document",
			idx: &index{
				dict: map[string][]int{},
			},
			r: strings.NewReader("Hello, world!"),
			wantDict: map[string][]int{
				"hello": {0},
				"world": {0},
			},
		},
		{
			name: "ok, second document",
			idx: &index{
				dict: map[string][]int{
					"hello": {0},
				},
				docIDCounter: 1,
			},
			r: strings.NewReader("Hello, world!"),
			wantDict: map[string][]int{
				"hello": {0, 1},
				"world": {1},
			},
		},
	}

	for _, tt := range tests {
		tt.idx.IndexDocument(tt.r)
		require.EqualValues(t, tt.wantDict, tt.idx.dict)
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
			err:   errTokenNotFound("x"),
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
		idx := &index{
			dict: tt.dict,
		}

		res, err := idx.GetPostingsList(tt.token)
		require.Equal(t, tt.err, err)
		require.Equal(t, tt.res, res)
	}
}
