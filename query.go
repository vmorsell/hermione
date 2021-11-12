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
	var res []int

	for _, t := range tokens {
		postings, err := q.idx.Postings(t)
		if err != nil {
			return nil, fmt.Errorf("get postings list: %w", err)
		}

		// Start with the full first postings list.
		if res == nil {
			res = postings
			fmt.Println(res)
			continue
		}
		res = intersect(res, postings)
		fmt.Println(res)
	}
	return res, nil
}

// intersect returns the intersection between two postings lists a and b.
func intersect(a, b []int) []int {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}

	res := []int{}

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
