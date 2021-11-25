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
				"x": {"a"},
			},
			tokens: []string{"x", "y"},
			err:    fmt.Errorf("get postings list: %w", errTokenNotInIndex("y")),
		},
		{
			name: "ok - two tokens",
			dict: map[string][]string{
				"x": {"a", "b", "c"},
				"y": {"b"},
			},
			tokens: []string{"x", "y"},
			calls: []intersectCall{
				{a: []string{"b"}, b: []string{"a", "b", "c"}},
			},
		},
		{
			name: "ok - three tokens",
			dict: map[string][]string{
				"x": {"a", "b", "c"},
				"y": {"b"},
				"z": {"b", "d"},
			},
			tokens: []string{"x", "y", "z"},
			calls: []intersectCall{
				{a: []string{"b"}, b: []string{"a", "b", "c"}}, // We expect it to start with the shortest list as 'a' param.
				{a: []string{"b"}, b: []string{"b", "d"}},
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
			b:    []string{"a"},
			res:  nil,
		},
		{
			name: "b postings list empty",
			a:    []string{"a"},
			b:    nil,
			res:  nil,
		},
		{
			name: "ok",
			a:    []string{"a", "b", "c"},
			b:    []string{"a", "c", "e"},
			res:  []string{"a", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := intersect(tt.a, tt.b)
			require.Equal(t, tt.res, res)
		})
	}
}
