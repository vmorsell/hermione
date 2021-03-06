package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntersection(t *testing.T) {
	type intersectCall struct {
		a, b []Posting
	}

	tests := []struct {
		name   string
		dict   map[string][]Posting
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
			dict: map[string][]Posting{
				"x": {
					{DocID: 0, Freq: 1},
				},
			},
			tokens: []string{"x", "y"},
			err:    fmt.Errorf("get postings list: %w", errTokenNotInIndex("y")),
		},
		{
			name: "ok - two tokens",
			dict: map[string][]Posting{
				"x": {
					{DocID: 0, Freq: 1},
					{DocID: 1, Freq: 1},
					{DocID: 2, Freq: 1},
				},
				"y": {
					{DocID: 1, Freq: 1},
				},
			},
			tokens: []string{"x", "y"},
			calls: []intersectCall{
				{
					a: []Posting{
						{DocID: 1, Freq: 1},
					},
					b: []Posting{
						{DocID: 0, Freq: 1},
						{DocID: 1, Freq: 1},
						{DocID: 2, Freq: 1},
					},
				},
			},
		},
		{
			name: "ok - three tokens",
			dict: map[string][]Posting{
				"x": {
					{DocID: 0, Freq: 1},
					{DocID: 1, Freq: 1},
					{DocID: 2, Freq: 1},
				},
				"y": {
					{DocID: 1, Freq: 1},
				},
				"z": {
					{DocID: 1, Freq: 1},
					{DocID: 3, Freq: 1},
				},
			},
			tokens: []string{"x", "y", "z"},
			calls: []intersectCall{
				{
					a: []Posting{
						{DocID: 1, Freq: 1}, // We expect it to start with the shortest list as 'a' param.
					},
					b: []Posting{
						{DocID: 0, Freq: 1},
						{DocID: 1, Freq: 1},
						{DocID: 2, Freq: 1},
					},
				},
				{
					a: []Posting{
						{DocID: 1, Freq: 1},
					},
					b: []Posting{
						{DocID: 1, Freq: 1},
						{DocID: 3, Freq: 1},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex().(*index)
			idx.dict = tt.dict

			q := NewQuerier(idx).(*querier)

			var calls []intersectCall
			q.intersectionFn = func(a, b []Posting) []Posting {
				calls = append(calls, intersectCall{a, b})
				return intersection(a, b)
			}

			_, err := q.Intersection(tt.tokens...)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.calls, calls)
		})
	}
}

func TestPrivateIntersection(t *testing.T) {
	tests := []struct {
		name string
		a    []Posting
		b    []Posting
		res  []Posting
	}{
		{
			name: "a postings list empty",
			a:    nil,
			b:    []Posting{{DocID: 0, Freq: 1}},
			res:  nil,
		},
		{
			name: "b postings list empty",
			a:    []Posting{{DocID: 0, Freq: 1}},
			b:    nil,
			res:  nil,
		},
		{
			name: "ok",
			a:    []Posting{{DocID: 0, Freq: 1}, {DocID: 1, Freq: 1}, {DocID: 2, Freq: 1}},
			b:    []Posting{{DocID: 0, Freq: 1}, {DocID: 2, Freq: 1}, {DocID: 4, Freq: 1}},
			res:  []Posting{{DocID: 0, Freq: 1}, {DocID: 2, Freq: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := intersection(tt.a, tt.b)
			require.Equal(t, tt.res, res)
		})
	}
}

func TestPrivatePhrase(t *testing.T) {
	tests := []struct {
		name string
		a, b []Posting
		res  []Posting
	}{
		{
			name: "no matches",
			a: []Posting{
				{
					DocID:     0,
					Positions: []int{0},
				},
				{
					DocID:     10,
					Positions: []int{10},
				},
			},
			b: []Posting{
				{
					DocID:     10,
					Positions: []int{1},
				},
			},
		},
		{
			name: "match",
			a: []Posting{
				{
					DocID:     10,
					Positions: []int{0, 4, 10},
				},
			},
			b: []Posting{
				{
					DocID:     0,
					Positions: []int{1, 5, 11},
				},
				{
					DocID:     10,
					Positions: []int{2, 5, 11, 20},
				},
			},
			res: []Posting{
				{
					DocID:     10,
					Positions: []int{5, 11},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQuerier(nil).(*querier)
			res := q.phrase(tt.a, tt.b)
			require.Equal(t, tt.res, res)
		})
	}
}
