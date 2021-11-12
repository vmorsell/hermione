package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntersectExported(t *testing.T) {
	tests := []struct {
		name   string
		q      *querier
		tokens []string
		res    []int
		err    error
	}{
		{
			name: "not ok - token not found",
			q: &querier{
				idx: &index{
					dict: map[string][]int{
						"a": {0},
					},
				},
			},
			tokens: []string{"a", "b"},
			err:    fmt.Errorf("get postings list: %w", errTokenNotInIndex("b")),
		},
		{
			name: "ok - two tokens",
			q: &querier{
				idx: &index{
					dict: map[string][]int{
						"a": {0, 1, 10, 100},
						"b": {1},
						"c": {1, 10, 100},
					},
				},
			},
			tokens: []string{"a", "c"},
			res:    []int{1, 10, 100},
		},
		{
			name: "ok - three tokens",
			q: &querier{
				idx: &index{
					dict: map[string][]int{
						"a": {0, 1, 10, 100},
						"b": {1},
						"c": {1, 10, 100},
					},
				},
			},
			tokens: []string{"a", "b", "c"},
			res:    []int{1},
		},
	}

	for _, tt := range tests {
		res, err := tt.q.Intersect(tt.tokens...)
		require.Equal(t, tt.err, err)
		require.Equal(t, tt.res, res)
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		res  []int
	}{
		{
			name: "a postings list empty",
			a:    nil,
			b:    []int{1},
			res:  nil,
		},
		{
			name: "b postings list empty",
			a:    []int{1},
			b:    nil,
			res:  nil,
		},
		{
			name: "ok",
			a:    []int{0, 10, 200, 1000},
			b:    []int{10, 250, 500, 1000, 1001},
			res:  []int{10, 1000},
		},
	}

	for _, tt := range tests {
		res := intersect(tt.a, tt.b)
		require.Equal(t, tt.res, res)
	}
}
