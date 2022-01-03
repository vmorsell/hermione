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
		dict         map[string][]Posting
		docIDCounter int
		idFn         func() (string, error)
		r            io.Reader
		wantDict     map[string][]Posting
		err          error
	}{
		{
			name: "ok - first document in index",
			dict: map[string][]Posting{},
			idFn: func() (string, error) { return "doc1", nil },
			r:    strings.NewReader("Hello, world!"),
			wantDict: map[string][]Posting{
				"hello": {
					{DocID: "doc1", Freq: 1},
				},
				"world": {
					{DocID: "doc1", Freq: 1},
				},
			},
		},
		{
			name: "ok - second document in index",
			dict: map[string][]Posting{
				"hello": {{DocID: "doc1", Freq: 1}},
			},
			idFn: func() (string, error) { return "doc2", nil },
			r:    strings.NewReader("Hello, world!"),
			wantDict: map[string][]Posting{
				"hello": {
					{DocID: "doc1", Freq: 1},
					{DocID: "doc2", Freq: 1},
				},
				"world": {
					{DocID: "doc2", Freq: 1},
				},
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
		dict  map[string][]Posting
		token string
		res   []Posting
		err   error
	}{
		{
			name: "not ok, token not found in dict",
			dict: map[string][]Posting{
				"hello": {
					{DocID: "doc1", Freq: 1},
				},
			},
			token: "world",
			err:   errTokenNotInIndex("world"),
		},
		{
			name: "ok",
			dict: map[string][]Posting{
				"hello": {
					{DocID: "doc1", Freq: 1},
				},
			},
			token: "hello",
			res: []Posting{
				{DocID: "doc1", Freq: 1},
			},
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

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		n    []int
		res  int
	}{
		{
			name: "no values",
			res:  0,
		},
		{
			name: "one value",
			n:    []int{1},
			res:  1,
		},
		{
			name: "multiple values",
			n:    []int{1, 10, 5},
			res:  10,
		},
	}

	for _, tt := range tests {
		res := max(tt.n...)
		require.Equal(t, tt.res, res)
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		n    []int
		res  int
	}{
		{
			name: "no values",
			res:  0,
		},
		{
			name: "one value",
			n:    []int{1},
			res:  1,
		},
		{
			name: "multiple values",
			n:    []int{1, 10, 5},
			res:  1,
		},
	}

	for _, tt := range tests {
		res := min(tt.n...)
		require.Equal(t, tt.res, res)
	}
}
