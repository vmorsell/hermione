package main

import "fmt"

type Querier interface {
	Intersect(tokens ...string) ([]int, error)
}

type querier struct {
	idx Index
}

func NewQuerier(idx Index) Querier {
	return &querier{
		idx: idx,
	}
}

func (q *querier) Intersect(tokens ...string) ([]int, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens provided")
	}

	postingsLists := make([][]int, 0, len(tokens))
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

		res = intersect(res, l)
	}
	return res, nil
}

// intersect returns the intersection between two postings lists a and b.
func intersect(a, b []int) []int {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}

	var res []int

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
