package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasDocID(t *testing.T) {
	tests := []struct {
		name         string
		postingsList []Posting
		id           int
		found        bool
		index        int
	}{
		{
			name: "found",
			postingsList: []Posting{
				{DocID: 0, Freq: 1},
			},
			id:    0,
			found: true,
			index: 0,
		},
		{
			name: "not found",
			postingsList: []Posting{
				{DocID: 0, Freq: 1},
			},
			id: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, index := hasDocID(tt.postingsList, tt.id)
			require.Equal(t, tt.found, found)
			require.Equal(t, tt.index, index)
		})
	}
}

func TestIndexDocument(t *testing.T) {
	tests := []struct {
		name     string
		dict     map[string][]Posting
		r        io.Reader
		wantDict map[string][]Posting
		err      error
	}{
		{
			name: "ok - first document in index",
			dict: map[string][]Posting{},
			r:    strings.NewReader("Hello hello, world!"),
			wantDict: map[string][]Posting{
				"hello": {
					{
						DocID:     0,
						Freq:      2,
						Positions: []int{0, 1},
					},
				},
				"world": {
					{
						DocID:     0,
						Freq:      1,
						Positions: []int{2},
					},
				},
			},
		},
		{
			name: "ok - second document in index",
			dict: map[string][]Posting{
				"hello": {
					{
						DocID:     0,
						Freq:      1,
						Positions: []int{0},
					},
				},
			},
			r: strings.NewReader("Hello hello, world!"),
			wantDict: map[string][]Posting{
				"hello": {
					{
						DocID:     0,
						Freq:      1,
						Positions: []int{0},
					},
					{
						DocID:     1,
						Freq:      2,
						Positions: []int{0, 1},
					},
				},
				"world": {
					{
						DocID:     1,
						Freq:      1,
						Positions: []int{2},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex().(*index)
			idx.dict = tt.dict
			idx.nextID = len(tt.dict)

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
					{DocID: 0, Freq: 1},
				},
			},
			token: "world",
			err:   errTokenNotInIndex("world"),
		},
		{
			name: "ok",
			dict: map[string][]Posting{
				"hello": {
					{DocID: 0, Freq: 1},
				},
			},
			token: "hello",
			res: []Posting{
				{DocID: 0, Freq: 1},
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

func TestID(t *testing.T) {
	index := NewIndex().(*index)

	for want := 0; want < 3; want++ {
		id := index.id()
		require.Equal(t, want, id)
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
