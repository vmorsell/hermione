package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexDocument(t *testing.T) {
	tests := []struct {
		name         string
		dict         map[string][]string
		docIDCounter int
		idFn         func() (string, error)
		r            io.Reader
		wantDict     map[string][]string
		err          error
	}{
		{
			name: "ok, first document",
			dict: map[string][]string{},
			idFn: func() (string, error) { return "a", nil },
			r:    strings.NewReader("Hello, world!"),
			wantDict: map[string][]string{
				"hello": {"a"},
				"world": {"a"},
			},
		},
		{
			name: "ok, second document",
			dict: map[string][]string{
				"hello": {"a"},
			},
			idFn: func() (string, error) { return "b", nil },
			r:    strings.NewReader("Hello, world!"),
			wantDict: map[string][]string{
				"hello": {"a", "b"},
				"world": {"b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex().(*index)
			idx.dict = tt.dict
			idx.idFn = tt.idFn

			idx.IndexDocument(tt.r)
			require.EqualValues(t, tt.wantDict, idx.dict)
		})
	}
}

func TestGetPostingsList(t *testing.T) {
	tests := []struct {
		name  string
		dict  map[string][]string
		token string
		res   []string
		err   error
	}{
		{
			name: "not ok, token not found in dict",
			dict: map[string][]string{
				"other": {"a"},
			},
			token: "x",
			err:   errTokenNotInIndex("x"),
		},
		{
			name: "ok",
			dict: map[string][]string{
				"x": {"a"},
			},
			token: "x",
			res:   []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex().(*index)
			idx.dict = tt.dict

			res, err := idx.Postings(tt.token)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.res, res)
		})
	}
}
