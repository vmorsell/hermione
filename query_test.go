package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntersect(t *testing.T) {
	type intersectCall struct {
		a, b []string
	}

	tests := []struct {
		name   string
		dict   map[string][]string
		tokens []string
		calls  []intersectCall
		err    error
	}{
		{
			name: "not ok - no tokens provided",
			err:  fmt.Errorf("no tokens provided"),
		},
		{
			name: "not ok - token not found",
			dict: map[string][]string{
				"x": {"doc1"},
			},
			tokens: []string{"x", "y"},
			err:    fmt.Errorf("get postings list: %w", errTokenNotInIndex("y")),
		},
		{
			name: "ok - two tokens",
			dict: map[string][]string{
				"x": {"doc1", "doc2", "doc3"},
				"y": {"doc2"},
			},
			tokens: []string{"x", "y"},
			calls: []intersectCall{
				{a: []string{"doc2"}, b: []string{"doc1", "doc2", "doc3"}},
			},
		},
		{
			name: "ok - three tokens",
			dict: map[string][]string{
				"x": {"doc1", "doc2", "doc3"},
				"y": {"doc2"},
				"z": {"doc2", "doc4"},
			},
			tokens: []string{"x", "y", "z"},
			calls: []intersectCall{
				{a: []string{"doc2"}, b: []string{"doc1", "doc2", "doc3"}}, // We expect it to start with the shortest list as 'a' param.
				{a: []string{"doc2"}, b: []string{"doc2", "doc4"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex().(*index)
			idx.dict = tt.dict

			q := NewQuerier(idx).(*querier)

			var calls []intersectCall
			q.intersectFn = func(a, b []string) []string {
				calls = append(calls, intersectCall{a, b})
				return intersect(a, b)
			}

			_, err := q.Intersect(tt.tokens...)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.calls, calls)
		})
	}
}

func TestPrivateIntersect(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		res  []string
	}{
		{
			name: "a postings list empty",
			a:    nil,
			b:    []string{"doc1"},
			res:  nil,
		},
		{
			name: "b postings list empty",
			a:    []string{"doc1"},
			b:    nil,
			res:  nil,
		},
		{
			name: "ok",
			a:    []string{"doc1", "doc2", "doc3"},
			b:    []string{"doc1", "doc3", "doc5"},
			res:  []string{"doc1", "doc3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := intersect(tt.a, tt.b)
			require.Equal(t, tt.res, res)
		})
	}
}
