package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

	type intersectCall struct {
		a, b []int
	}

	tests := []struct {
		name   string
		dict   map[string][]int
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
			dict: map[string][]int{
				"a": {0},
			},
			tokens: []string{"a", "b"},
			err:    fmt.Errorf("get postings list: %w", errTokenNotInIndex("b")),
		},
		{
			name: "ok - two tokens",
			dict: map[string][]int{
				"a": {0, 1, 2},
				"b": {1, 3},
			},
			tokens: []string{"a", "b"},
			calls: []intersectCall{
				{a: []int{1, 3}, b: []int{0, 1, 2}},
			},
		},
		{
			name: "ok - three tokens",
			dict: map[string][]int{
				"a": {0, 1, 2},
				"b": {1, 3},
				"c": {1},
			},
			tokens: []string{"a", "b", "c"},
			calls: []intersectCall{
				{a: []int{1}, b: []int{0, 1, 2}}, // We expect it to start with the shortest list as 'a' param.
				{a: []int{1}, b: []int{1, 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex().(*index)
			idx.dict = tt.dict

			q := NewQuerier(idx).(*querier)

			var calls []intersectCall
			q.intersectFn = func(a, b []int) []int {
				calls = append(calls, intersectCall{a, b})
				return intersect(a, b)
			}

			_, err := q.Intersect(tt.tokens...)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.calls, calls)
		})
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
