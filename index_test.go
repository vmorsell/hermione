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

