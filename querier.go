package main

import (
	"fmt"
	"strings"
)

type Querier interface {
	Boolean(query string) ([]string, error)
	Intersect(tokens ...string) ([]string, error)
}

type querier struct {
	idx Index

	intersectFn func(a, b []string) []string
}

func NewQuerier(idx Index) Querier {
	return &querier{
		idx: idx,

		// Use the default functions
		intersectFn: intersect,
	}
}

// Boolean takes a query expression and returns matching documents using
// boolean retrieval.
func (q *querier) Boolean(query string) ([]string, error) {
	tokens := strings.Split(query, " ")
	return q.Intersect(tokens...)
}

// Intersect fetches the postings lists for all given terms and returns the
// document ID's present in all lists.
func (q *querier) Intersect(tokens ...string) ([]string, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens provided")
	}

	postingsLists := make([][]string, 0, len(tokens))
	var lowestDocFreqIdx int

	// Fetch all postings lists.
	for _, t := range tokens {
		postings, err := q.idx.Postings(t)
		if err != nil {
			return nil, fmt.Errorf("get postings list: %w", err)
		}
		postingsLists = append(postingsLists, postings)

		// Keep track of the shortest postings list. We will start with
		// that one when doing the intersection.
		if len(postingsLists[lowestDocFreqIdx]) > len(postings) {
			lowestDocFreqIdx = len(postingsLists) - 1
		}
	}

	res := postingsLists[lowestDocFreqIdx]
	for i, l := range postingsLists {
		if i == lowestDocFreqIdx {
			continue
		}

		res = q.intersectFn(res, l)
	}
	return res, nil
}

// intersect returns the common document ID's from the two given postings lists.
func intersect(a, b []string) []string {
	if len(a) == 0 {
		return nil
	}
	if len(b) == 0 {
		return nil
	}

	var res []string

	aCur := 0
	bCur := 0

	for {
		if aCur > len(a)-1 || bCur > len(b)-1 {
			break
		}

		if a[aCur] == b[bCur] {
			res = append(res, a[aCur])
			aCur++
			bCur++
			continue
		}

		if a[aCur] > b[bCur] {
			bCur++
		} else {
			aCur++
		}
	}
	return res
}
